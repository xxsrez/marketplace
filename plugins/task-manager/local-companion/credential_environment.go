package main

import (
	"context"
	"strings"
	"sync"
)

const localTokenEnvironment = "TASK_MANAGER_LOCAL_TOKEN"

// Linux uses an explicitly provisioned REST credential, never a desktop OAuth
// cache or a browser callback. The host owns secret injection and rotation.
type environmentTokenSource struct {
	mutex    sync.Mutex
	lookup   func(string) string
	rejected bool
}

func (source *environmentTokenSource) Token(ctx context.Context) (string, error) {
	source.mutex.Lock()
	defer source.mutex.Unlock()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if source.rejected {
		return "", newLocalError("environment_token_rejected", "Task Manager rejected TASK_MANAGER_LOCAL_TOKEN; replace the environment secret and restart the companion")
	}
	token := source.lookup(localTokenEnvironment)
	if token == "" {
		return "", newLocalError("environment_token_required", "Linux local uploads require TASK_MANAGER_LOCAL_TOKEN: inject a Task Manager personal API token with api:write through the host secret environment, then restart the companion; do not paste it into chat or tool arguments")
	}
	if !strings.HasPrefix(token, "tm_pat_") || len(token) <= len("tm_pat_") || len(token) > 4096 ||
		strings.ContainsFunc(token, func(r rune) bool {
			return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-')
		}) {
		return "", newLocalError("environment_token_invalid", "TASK_MANAGER_LOCAL_TOKEN must be a Task Manager personal API token without whitespace; configure it in the host secret environment")
	}
	return token, nil
}

func (source *environmentTokenSource) Invalidate() {
	source.mutex.Lock()
	defer source.mutex.Unlock()
	source.rejected = true
}
