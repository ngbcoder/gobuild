package runner

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// TaskRunner owns the exec lifecycle for a single task.
type TaskRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

// NewTaskRunner creates a cancellable runner for a shell-backed task.
func NewTaskRunner(parent context.Context) *TaskRunner {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &TaskRunner{ctx: ctx, cancel: cancel}
}

// Context returns the runner-scoped context.
func (r *TaskRunner) Context() context.Context {
	return r.ctx
}

// Cancel stops any active task execution.
func (r *TaskRunner) Cancel() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
}

// Run executes a command using the runner-scoped context and returns the combined output.
func (r *TaskRunner) Run(name string, args ...string) (string, error) {
	cmd := exec.CommandContext(r.ctx, name, args...)
	return captureOutput(cmd)
}

func captureOutput(cmd *exec.Cmd) (string, error) {
	if cmd == nil {
		return "", nil
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		combined := stdout.String() + stderr.String()
		if strings.TrimSpace(combined) == "" {
			return "", err
		}
		return combined, err
	}

	combined := stdout.String() + stderr.String()
	return strings.TrimSpace(combined), nil
}

// RunWithTimeout executes a command and terminates it if it exceeds the timeout.
func RunWithTimeout(parent context.Context, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	return captureOutput(cmd)
}
