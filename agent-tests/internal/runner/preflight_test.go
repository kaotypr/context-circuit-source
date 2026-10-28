package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
)

func TestRunStopsWhenAuthIsMissing(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "started")
	stub := writeTouchStub(t, dir, marker)
	hostsFile := writeHostsFile(t, dir, []hostRow{
		{name: "claude", command: stub, model: "claude-opus-5-5", version: "0", timeout: 5, auth: []string{"ANTHROPIC_API_KEY"}},
		{name: "cursor", command: stub, model: "composer-2.5", version: "0", timeout: 5, auth: []string{"CURSOR_API_KEY"}},
	})
	scenarios := writeScenario(t, dir, "marker")
	out := filepath.Join(dir, "out")
	var buf bytes.Buffer
	err := Run(context.Background(), RunOptions{
		Checkout:     dir,
		Runs:         1,
		HostsFile:    hostsFile,
		ScenariosDir: scenarios,
		OutDir:       out,
		Caller:       map[string]string{"PATH": "/usr/bin:/bin", "ANTHROPIC_API_KEY": "stub-anthropic-auth-value"},
		Confirm:      true,
		Stdout:       &buf,
	})
	if err == nil || !strings.Contains(err.Error(), "CURSOR_API_KEY") || !strings.Contains(err.Error(), "no host was started") {
		t.Fatalf("err %v\n%s", err, buf.String())
	}
	if strings.Contains(err.Error(), "claude") {
		t.Fatalf("a ready host is named as a problem: %v", err)
	}
	if strings.Contains(buf.String(), "stub-anthropic-auth-value") {
		t.Fatal("output contains an auth value")
	}
	if !strings.Contains(buf.String(), "auth claude: ready") || !strings.Contains(buf.String(), "auth cursor: not ready (missing CURSOR_API_KEY)") {
		t.Fatalf("readiness lines\n%s", buf.String())
	}
	if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
		t.Fatal("a host started without auth")
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Fatal("output directory was created before the auth check")
	}
}

func TestDryRunReportsAuthWithoutStopping(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "started")
	stub := writeTouchStub(t, dir, marker)
	hostsFile := writeHostsFile(t, dir, []hostRow{
		{name: "cursor", command: stub, model: "composer-2.5", version: "0", timeout: 5, auth: []string{"CURSOR_API_KEY"}},
	})
	scenarios := writeScenario(t, dir, "marker")
	var buf bytes.Buffer
	err := Run(context.Background(), RunOptions{
		Runs:         1,
		HostsFile:    hostsFile,
		ScenariosDir: scenarios,
		Caller:       map[string]string{"PATH": "/usr/bin:/bin"},
		DryRun:       true,
		Stdout:       &buf,
	})
	if err != nil {
		t.Fatalf("dry run: %v\n%s", err, buf.String())
	}
	if !strings.Contains(buf.String(), "auth cursor: not ready (missing CURSOR_API_KEY)") {
		t.Fatalf("dry run did not report readiness\n%s", buf.String())
	}
	if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
		t.Fatal("dry run started a host")
	}
}

func TestMeasureHostVersions(t *testing.T) {
	dir := t.TempDir()
	good := writeVersionStub(t, dir, filepath.Join(dir, "started"))
	bad := filepath.Join(dir, "bad")
	if err := os.WriteFile(bad, []byte("#!/bin/sh\nexit 3\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	got := measureHostVersions(ctx, []hosts.HostConfig{
		{Name: "claude", Command: good, HostVersion: "configured-1"},
		{Name: "cursor", Command: bad, HostVersion: "configured-2"},
	}, "/usr/bin:/bin")
	if hv := got["claude"]; hv.value != "stub-9" || !hv.measured {
		t.Fatalf("claude %+v", hv)
	}
	if hv := got["cursor"]; hv.value != "configured-2" || hv.measured {
		t.Fatalf("cursor %+v", hv)
	}
}

func TestRecordLineHostVersionSource(t *testing.T) {
	item := planned{
		Scenario: scenario.Scenario{ID: "marker", Version: 1},
		Host:     "cursor",
		Config:   hosts.HostConfig{Name: "cursor", HostVersion: "configured-2"},
		Side:     "branch",
		RunN:     1,
	}
	measured := &runtime{hostVersions: map[string]hostVersion{"cursor": {value: "real-3", measured: true}}}
	line := measured.recordLine(item, hosts.Parsed{}, "pass", nil, 0, "", nil, nil)
	if line.HostVersion != "real-3" || line.HostVersionSource != "command" || line.NullReasons["host_version"] != "" {
		t.Fatalf("measured %+v", line)
	}
	fallback := &runtime{hostVersions: map[string]hostVersion{"cursor": {value: "configured-2"}}}
	line = fallback.recordLine(item, hosts.Parsed{}, "pass", nil, 0, "", nil, nil)
	if line.HostVersion != "configured-2" || line.HostVersionSource != "config" || line.NullReasons["host_version"] == "" {
		t.Fatalf("fallback %+v", line)
	}
}
