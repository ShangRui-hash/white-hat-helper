package memory

import (
	"context"
	"errors"
)

var taskCancelFuncMap map[int64]context.CancelFunc

func init() {
	taskCancelFuncMap = make(map[int64]context.CancelFunc)
}

func RegisterTaskCancelFunc(taskID int64, f context.CancelFunc) {
	taskCancelFuncMap[taskID] = f
}

func StopTask(taskID int64) error {
	f, ok := taskCancelFuncMap[taskID]
	if !ok {
		return errors.New("no cancel func")
	}
	f()
	return nil
}
