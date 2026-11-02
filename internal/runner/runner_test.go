package runner

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRunnerCancelsContext(t *testing.T) {
	r := NewTaskRunner(context.Background())
	if r == nil {
		t.Fatal("NewTaskRunner() returned nil")
	}
	if r.Context() == nil {
		t.Fatal("Context() returned nil")
	}

	r.Cancel()
	select {
	case <-r.Context().Done():
		// expected cancelation signal
	default:
		t.Fatal("Context() was not canceled after Cancel()")
	}
}

func TestExecCommandCapturesOutput(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'hello\\nworld\\n'")
	out, err := captureOutput(cmd)
	if err != nil {
		t.Fatalf("captureOutput() returned error: %v", err)
	}
	if !strings.Contains(out, "hello") || !strings.Contains(out, "world") {
		t.Fatalf("output = %q, want both lines captured", out)
	}
}

func TestExecCommandHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", "-c", "sleep 5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	cancel()
	if err := cmd.Wait(); err == nil {
		t.Fatal("Wait() succeeded despite cancelation, want non-nil error")
	}
}

func TestCaptureOutputTimeout(t *testing.T) {
	cmd := exec.Command("sh", "-c", "sleep 0.25; printf 'done' ")
	start := time.Now()
	out, err := captureOutput(cmd)
	if err != nil {
		t.Fatalf("captureOutput() returned error: %v", err)
	}
	if out != "done" {
		t.Fatalf("output = %q, want %q", out, "done")
	}
	if time.Since(start) > time.Second {
		t.Fatal("captureOutput() took too long")
	}
}
