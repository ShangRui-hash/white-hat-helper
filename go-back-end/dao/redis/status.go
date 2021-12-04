package redis

import "go.uber.org/zap"

const (
	NOT_STARTED_TASK = "not_started_task"
	RUNNING_TASK     = "running_task"
	STOPPED_TASK     = "stopped_task"
	COMPLETED_TASK   = "completed_task"
	FAILED_TASK      = "failed_task"
)

var statusMap = map[string]string{
	NOT_STARTED_TASK: "未开始",
	RUNNING_TASK:     "运行中",
	STOPPED_TASK:     "暂停",
	COMPLETED_TASK:   "完成",
	FAILED_TASK:      "失败",
}

//GetTaskStatus 获取任务状态
func GetTaskStatus(taskID int64) (status string, err error) {
	for key := range statusMap {
		ok, err := rdb.SIsMember(key, taskID).Result()
		if err != nil {
			return "", err
		}
		if ok {
			return key, nil
		}
	}
	return "", nil
}
func GetTaskStatusText(taskID int64) (status string, err error) {
	status, err = GetTaskStatus(taskID)
	if err != nil {
		return "", err
	}
	return statusMap[status], nil
}

//SetTaskStatus 设置任务的状态
func SetTaskStatus(taskID int64, statusKey string) (err error) {
	//1.从原集合中移除
	if err := DeleteTaskStatus(taskID); err != nil {
		zap.L().Error("从原集合中移除任务失败", zap.Error(err))
		return err
	}
	//2.加入到新集合中
	_, err = rdb.SAdd(statusKey, taskID).Result()
	return err
}

//InitTaskStatus 初始化任务状态
func InitTaskStatus(taskID int64) (err error) {
	return SetTaskStatus(taskID, NOT_STARTED_TASK)
}

func BeginTask(taskID int64) (err error) {
	return SetTaskStatus(taskID, RUNNING_TASK)
}

func StopTask(taskID int64) (err error) {
	return SetTaskStatus(taskID, STOPPED_TASK)
}

func CompletedTask(taskID int64) (err error) {
	return SetTaskStatus(taskID, COMPLETED_TASK)
}

func DeleteTaskStatus(taskID int64) (err error) {
	oldStatusKey, err := GetTaskStatus(taskID)
	if err != nil {
		return err
	}
	if oldStatusKey != "" {
		_, err := rdb.SRem(oldStatusKey, taskID).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
