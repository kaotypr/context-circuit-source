package proc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func skipShell(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell fixtures")
	}
}

func writeScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cmd.sh")
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestEmptyArgv(t *testing.T) {
	_, err := Run(context.Background(), Spec{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStartError(t *testing.T) {
	_, err := Run(context.Background(), Spec{Argv: []string{"/no/such/proc-binary"}})
	if err == nil {
		t.Fatal("expected start error")
	}
}

func TestExitCode(t *testing.T) {
	skipShell(t)
	script := writeScript(t, "#!/bin/sh\nexit 17\n")
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/usr/bin:/bin"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 17 {
		t.Fatalf("exit %d", res.ExitCode)
	}
	if res.TimedOut {
		t.Fatal("timed out")
	}
}

func TestZeroExitIsNilError(t *testing.T) {
	skipShell(t)
	script := writeScript(t, "#!/bin/sh\nexit 0\n")
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/usr/bin:/bin"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit %d", res.ExitCode)
	}
}

func TestEnvIsNotInherited(t *testing.T) {
	skipShell(t)
	const sentinel = "PROC_SENTINEL_SHOULD_NOT_APPEAR"
	t.Setenv(sentinel, "leaked")
	script := writeScript(t, "#!/bin/sh\nexec /usr/bin/env\n")
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/usr/bin:/bin", "ONLY=yes"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	out := string(res.Stdout)
	if !strings.Contains(out, "ONLY=yes") {
		t.Fatalf("missing ONLY:\n%s", out)
	}
	if strings.Contains(out, sentinel) {
		t.Fatalf("inherited sentinel:\n%s", out)
	}
	if home := os.Getenv("HOME"); home != "" && strings.Contains(out, "HOME="+home) {
		t.Fatalf("inherited HOME:\n%s", out)
	}
}

func TestOnLineOrder(t *testing.T) {
	skipShell(t)
	script := writeScript(t, "#!/bin/sh\nprintf 'one\\n'\nprintf 'two\\n'\nprintf 'three\\n'\n")
	var got []string
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/usr/bin:/bin"},
		Timeout: 10 * time.Second,
		OnLine: func(line []byte) {
			got = append(got, string(line))
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"one", "two", "three"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("lines %q stdout %q", got, res.Stdout)
	}
}

func TestLongLine(t *testing.T) {
	skipShell(t)
	const n = 5 << 20
	script := writeScript(t, "#!/bin/sh\n/bin/dd if=/dev/zero bs="+itoa(n)+" count=1 2>/dev/null | /usr/bin/tr '\\0' 'A'\nprintf '\\n'\n")
	var lines []int
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/usr/bin:/bin"},
		Timeout: 20 * time.Second,
		OnLine: func(line []byte) {
			lines = append(lines, len(line))
			if len(line) > 0 && line[0] != 'A' {
				t.Errorf("line starts with %q", line[0])
			}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != n {
		t.Fatalf("online lens %v", lines)
	}
	if len(res.Stdout) != n+1 || res.Stdout[len(res.Stdout)-1] != '\n' {
		t.Fatalf("stdout len %d", len(res.Stdout))
	}
	for i, b := range res.Stdout[:n] {
		if b != 'A' {
			t.Fatalf("stdout byte %d is %q", i, b)
		}
	}
}

func TestStreamTee(t *testing.T) {
	skipShell(t)
	dir := t.TempDir()
	script := filepath.Join(dir, "cmd.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf 'hello\\n'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	stream := filepath.Join(dir, "stream.txt")
	res, err := Run(context.Background(), Spec{
		Argv:       []string{script},
		Env:        []string{"PATH=/usr/bin:/bin"},
		Timeout:    10 * time.Second,
		StreamPath: stream,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(res.Stdout) != "hello\n" {
		t.Fatalf("stdout %q", res.Stdout)
	}
	got, err := os.ReadFile(stream)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\n" {
		t.Fatalf("stream %q", got)
	}
}

func TestTimeoutKillsGrandchild(t *testing.T) {
	skipShell(t)
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "grand.pid")
	body := fmt.Sprintf("#!/bin/sh\n/bin/sh -c 'echo $$ > %s; exec sleep 60' &\nwait\n", strconv.Quote(pidFile))
	script := filepath.Join(dir, "cmd.sh")
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	var pid int
	t.Cleanup(func() {
		if pid > 1 {
			killPid(pid)
		}
	})
	res, err := Run(context.Background(), Spec{
		Argv:    []string{script},
		Env:     []string{"PATH=/bin:/usr/bin"},
		Timeout: 1500 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.TimedOut {
		t.Fatalf("expected timeout, exit %d stderr %q", res.ExitCode, res.Stderr)
	}
	raw, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err = strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil || pid <= 1 {
		t.Fatalf("pid %q: %v", raw, err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("grandchild %d still alive", pid)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestCappedMarker(t *testing.T) {
	var c capped
	if _, err := c.Write([]byte("hi")); err != nil {
		t.Fatal(err)
	}
	if string(c.bytes()) != "hi" {
		t.Fatalf("small write %q", c.bytes())
	}

	var big capped
	chunk := bytes.Repeat([]byte("a"), 4096)
	written := 0
	for written < maxCapture+len(chunk) {
		n, err := big.Write(chunk)
		if err != nil || n != len(chunk) {
			t.Fatalf("write %d %v", n, err)
		}
		written += n
	}
	got := big.bytes()
	if len(got) != maxCapture {
		t.Fatalf("len %d want %d", len(got), maxCapture)
	}
	if !bytes.HasSuffix(got, []byte(truncMarker)) {
		t.Fatal("missing truncation marker")
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
