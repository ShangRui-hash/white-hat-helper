package redis

const (
	NOT_STARTED = "not started"
	RUNNING     = "running"
	STOPPED     = "stopped"
	COMPLETED   = "completed"
	FAILED      = "failed"
)

//GetTaskStatus 获取任务状态
func GetTaskStatus(taskID int64) (status map[string]string, err error) {
	return rdb.HGetAll(GetTaskStatusHashKey(taskID)).Result()
}

//SetTaskStatus 设置任务的某个工具状态
func SetTaskStatus(taskID int64, toolName, statusKey string) (err error) {
	return rdb.HSet(GetTaskStatusHashKey(taskID), toolName, statusKey).Err()
}

func BeginTask(taskID int64, toolName string) (err error) {
	return SetTaskStatus(taskID, toolName, RUNNING)
}

func StopTask(taskID int64, toolName string) (err error) {
	return SetTaskStatus(taskID, toolName, STOPPED)
}

func CompletedTask(taskID int64, toolName string) (err error) {
	return SetTaskStatus(taskID, toolName, COMPLETED)
}

func FailedTask(taskID int64, toolName string) (err error) {
	return SetTaskStatus(taskID, toolName, FAILED)
}
func DeleteTaskStatus(taskID int64) error {
	return rdb.Del(GetTaskStatusHashKey(taskID)).Err()
}
