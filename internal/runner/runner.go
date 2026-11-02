package runner

import "context"

// TaskRunner owns the exec lifecycle for a single task.
type TaskRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTaskRunner creates a cancellable runner for a shell-backed task.
func NewTaskRunner(parent context.Context) *TaskRunner {
	ctx, cancel := context.WithCancel(parent)
	return &TaskRunner{ctx: ctx, cancel: cancel}
}

// Context returns the runner-scoped context.
func (r *TaskRunner) Context() context.Context {
	return r.ctx
}

// Cancel stops any active task execution.
func (r *TaskRunner) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}
