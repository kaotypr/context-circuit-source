// Package proc runs a child process in its own process group.
package proc

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	// maxCapture is the in-memory cap for stdout and stderr.
	maxCapture  = 64 << 20
	truncMarker = "[truncated]"
	killGrace   = 5 * time.Second
)

// Spec describes one child process. Env is the complete environment:
// the child inherits nothing from this process.
type Spec struct {
	Argv []string
	Dir  string
	Env  []string
	// Stdin is copied to the child. An empty string supplies empty stdin
	// (the same bytes the child would read from /dev/null).
	Stdin      string
	Timeout    time.Duration
	StreamPath string
	// OnLine, when set, is called once per stdout line in order.
	// The slice is owned by the caller and does not include the newline.
	OnLine func([]byte)
}

// Result is the outcome of a child that was started.
// A non-zero ExitCode is not an error; failure to start is.
type Result struct {
	ExitCode int
	TimedOut bool
	Wall     time.Duration
	Stdout   []byte
	Stderr   []byte
}

// Run starts s and waits for it. On timeout or ctx cancellation the whole
// process group is signalled SIGTERM, then SIGKILL after a grace period.
func Run(ctx context.Context, s Spec) (Result, error) {
	if len(s.Argv) == 0 || s.Argv[0] == "" {
		return Result{}, errors.New("empty argv")
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	cmd := exec.Command(s.Argv[0], s.Argv[1:]...)
	cmd.Dir = s.Dir
	if s.Env == nil {
		cmd.Env = []string{}
	} else {
		cmd.Env = s.Env
	}
	cmd.Stdin = strings.NewReader(s.Stdin)
	prepareCmd(cmd)

	// Own pipes, not Cmd.StdoutPipe: Wait closes pipes it created, and that
	// races a concurrent reader. The write ends are closed in this process
	// after Start so the readers see EOF when the child exits.
	stdout, stdoutW, err := os.Pipe()
	if err != nil {
		return Result{}, err
	}
	stderr, stderrW, err := os.Pipe()
	if err != nil {
		stdout.Close()
		stdoutW.Close()
		return Result{}, err
	}
	cmd.Stdout = stdoutW
	cmd.Stderr = stderrW

	// A nil *os.File passed as an io.Writer is a non-nil interface and
	// Write returns ErrInvalid, so the tee stays a nil interface when unset.
	var tee io.Writer
	if s.StreamPath != "" {
		stream, err := os.Create(s.StreamPath)
		if err != nil {
			stdout.Close()
			stdoutW.Close()
			stderr.Close()
			stderrW.Close()
			return Result{}, err
		}
		defer stream.Close()
		tee = stream
	}

	started := time.Now()
	if err := cmd.Start(); err != nil {
		stdout.Close()
		stdoutW.Close()
		stderr.Close()
		stderrW.Close()
		return Result{}, err
	}
	stdoutW.Close()
	stderrW.Close()
	defer stdout.Close()
	defer stderr.Close()
	proc := cmd.Process

	var outBuf, errBuf capped
	var outErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		outErr = readStdout(stdout, &outBuf, tee, s.OnLine)
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&errBuf, stderr)
	}()

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	var timeout <-chan time.Time
	if s.Timeout > 0 {
		timer := time.NewTimer(s.Timeout)
		defer timer.Stop()
		timeout = timer.C
	}

	var waitErr error
	timedOut := false
	select {
	case waitErr = <-waitCh:
	case <-ctx.Done():
		waitErr = killAndWait(proc, waitCh)
	case <-timeout:
		timedOut = true
		waitErr = killAndWait(proc, waitCh)
	}
	wg.Wait()

	code, statusErr := exitCode(waitErr)
	res := Result{
		ExitCode: code,
		TimedOut: timedOut,
		Wall:     time.Since(started),
		Stdout:   outBuf.bytes(),
		Stderr:   errBuf.bytes(),
	}
	if statusErr != nil {
		return res, statusErr
	}
	if outErr != nil {
		return res, outErr
	}
	return res, nil
}

func killAndWait(proc *os.Process, waitCh <-chan error) error {
	signal(proc, false)
	select {
	case err := <-waitCh:
		return err
	case <-time.After(killGrace):
		signal(proc, true)
		return <-waitCh
	}
}

func exitCode(err error) (int, error) {
	if err == nil {
		return 0, nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), nil
	}
	return 0, err
}

func readStdout(r io.Reader, mem *capped, stream io.Writer, onLine func([]byte)) error {
	br := bufio.NewReader(r)
	for {
		line, err := readLine(br)
		if len(line) > 0 {
			_, _ = mem.Write(line)
			if stream != nil {
				if _, werr := stream.Write(line); werr != nil {
					return werr
				}
			}
			if onLine != nil {
				payload := line
				if n := len(payload); n > 0 && payload[n-1] == '\n' {
					payload = payload[:n-1]
					if n := len(payload); n > 0 && payload[n-1] == '\r' {
						payload = payload[:n-1]
					}
				}
				cp := make([]byte, len(payload))
				copy(cp, payload)
				onLine(cp)
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// readLine reads one line with no length limit. The returned slice includes
// the trailing newline when the line had one.
func readLine(br *bufio.Reader) ([]byte, error) {
	var buf []byte
	for {
		frag, err := br.ReadSlice('\n')
		buf = append(buf, frag...)
		if err == bufio.ErrBufferFull {
			continue
		}
		return buf, err
	}
}

// capped stores at most maxCapture bytes. Past that it keeps a prefix and
// a trailing [truncated] marker.
type capped struct {
	buf       []byte
	truncated bool
}

func (c *capped) Write(p []byte) (int, error) {
	if c.truncated || len(p) == 0 {
		return len(p), nil
	}
	room := maxCapture - len(truncMarker)
	if len(c.buf) >= room {
		c.buf = append(c.buf[:room], truncMarker...)
		c.truncated = true
		return len(p), nil
	}
	space := room - len(c.buf)
	if len(p) <= space {
		c.buf = append(c.buf, p...)
		return len(p), nil
	}
	c.buf = append(c.buf, p[:space]...)
	c.buf = append(c.buf, truncMarker...)
	c.truncated = true
	return len(p), nil
}

func (c *capped) bytes() []byte {
	return c.buf
}
