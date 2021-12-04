package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

func AddTask(param *models.ParamAddTask) (*models.Task, error) {
	//1.基本信息进入mysql
	id, err := mysql.AddTask(param)
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
	//2.获取运行状态
	status, err := redis.GetTaskStatusText(id)
	if err != nil {
		return nil, err
	}
	task.Status = status
	return task, nil
}

func GetTaskList(param *models.Page) ([]*models.Task, error) {
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
	}
	return tasks, nil
}

// func UpdateTask(param *models.ParamUpdateTask) error {
// 	//1.更新mysql
// 	err := mysql.UpdateTask(param)
// 	if err != nil {
// 		return err
// 	}
// 	//2.更新redis
// 	if err := redis.UpdateTaskStatus(param.ID); err != nil {
// 		return err
// 	}
// 	return nil
// }

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
