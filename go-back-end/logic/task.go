package logic

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/param"

	"go.uber.org/zap"
)

func AddTask(param *param.ParamAddTask) (*models.Task, error) {
	//1.基本信息进入mysql
	id, err := mysql.AddTask(param.Name, param.ScanArea, int64(param.CompanyID))
	if err != nil {
		zap.L().Error("mysql.AddTask failed", zap.Error(err))
		return nil, err
	}
	//2.频繁改动的信息存入redis
	if err := redis.InitTaskStatus(id); err != nil {
		zap.L().Error("redis.InitTaskStatus failed", zap.Error(err))
		return nil, err
	}
	//3.获取存入的信息，返回给前端
	return GetTaskByID(id)
}

func UpdateTask(param *param.ParamUpdateTask) (*models.Task, error) {
	//更新mysql
	err := mysql.UpdateTask(param.ID, param.CompanyID, param.Name, param.ScanArea)
	if err != nil {
		return nil, err
	}
	//获取存入的信息，返回给前端
	return GetTaskByID(param.ID)
}

func GetTaskByID(id int64) (*models.Task, error) {
	//1.获取基本信息
	task, err := mysql.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	company, err := mysql.GetCompanyByID(task.CompanyID)
	if err != nil {
		return nil, err
	}
	task.CompanyName = company.Name
	task.ScanAreaList = strings.Split(task.ScanArea, ",")
	//2.获取运行状态
	status, err := redis.GetTaskStatusText(id)
	if err != nil {
		return nil, err
	}
	task.Status = status
	return task, nil
}

func GetTaskList(param *param.Page) ([]*models.Task, error) {
	//1.获取基本信息
	tasks, err := mysql.GetTaskList(param.Offset, param.Count)
	if err != nil {
		return nil, err
	}
	//2.获取运行状态
	for _, task := range tasks {
		status, err := redis.GetTaskStatusText(task.ID)
		if err != nil {
			return nil, err
		}
		task.Status = status
		task.ScanAreaList = strings.Split(task.ScanArea, ",")
	}
	return tasks, nil
}

func DeleteTask(id int64) error {
	//1.删除mysql
	err := mysql.DeleteTask(id)
	if err != nil {
		return err
	}
	//2.删除redis
	if err := redis.DeleteTaskStatus(id); err != nil {
		return err
	}
	return nil
}

func StartTask(taskID int64) error {
	//1.查询任务的详细信息
	task, err := mysql.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	//2.启动任务
	logFile, err := os.OpenFile("task.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"white-hat-helper",
		"-r", "/Applications/MAMP/htdocs/安全开发/PythonProject/white-hat-helper/scanner/config/redis.json",
		"--cid", fmt.Sprintf("%d", task.CompanyID),
		"--dict", "/Applications/MAMP/htdocs/安全开发/PythonProject/white-hat-helper/scanner/dirsearch.txt",
		"-d", task.ScanArea)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		return err
	}
	//2.存储对应的进程pid
	if err := redis.SaveTaskPid(taskID, cmd.Process.Pid); err != nil {
		return err
	}
	//3.更新redis
	if err := redis.BeginTask(taskID); err != nil {
		return err
	}
	return nil
}

func StopTask(taskID int64) error {
	//1.获取任务对应的pid
	pid, err := redis.GetTaskPid(taskID)
	if err != nil {
		return err
	}
	//2.结束任务
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := process.Kill(); err != nil {
		return err
	}
	//3.删除redis中的pid
	if err := redis.ClearTaskPid(taskID); err != nil {
		return err
	}
	if err := redis.StopTask(taskID); err != nil {
		return err
	}
	return nil
}
