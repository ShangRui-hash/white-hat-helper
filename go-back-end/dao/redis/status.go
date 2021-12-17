package redis

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

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
	STOPPED_TASK:     "停止",
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

//GetTaskStatusText 获取任务状态的文本
func GetTaskStatusText(taskID int64) (status string, err error) {
	status, err = GetTaskStatus(taskID)
	if err != nil {
		return "", err
	}
	return statusMap[status], nil
}

//SaveTaskPid 存储任务id对应的pid
func SaveTaskPid(taskID int64, pid int) (err error) {
	return rdb.HSet(TaskPidHashKey, fmt.Sprintf("%v", taskID), pid).Err()
}

//ClearTaskStatus 清除任务对应的pid
func ClearTaskPid(taskID int64) (err error) {
	return rdb.HDel(TaskPidHashKey, fmt.Sprintf("%v", taskID)).Err()
}

//GetTaskPid 获取任务的pid
func GetTaskPid(taskID int64) (pid int, err error) {
	temp, err := rdb.HGet(TaskPidHashKey, fmt.Sprintf("%v", taskID)).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(temp)
}

//GetRunningTaskList 获取正在运行的任务列表
func GetRunningTaskList() (taskIDList []int64, err error) {
	taskIDList = make([]int64, 0)
	temp, err := rdb.SMembers(RUNNING_TASK).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range temp {
		taskID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		taskIDList = append(taskIDList, taskID)
	}
	return taskIDList, nil
}

func StopAllRunningTask() (err error) {
	taskIDList, err := GetRunningTaskList()
	if err != nil {
		return err
	}
	for _, taskID := range taskIDList {
		err = StopTask(taskID)
		if err != nil {
			zap.L().Error("StopAllRunningTask", zap.Error(err))
			return err
		}
	}
	return nil
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
