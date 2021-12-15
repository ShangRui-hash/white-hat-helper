package mysql

import (
	"strings"
	"web_app/models"

	"go.uber.org/zap"
)

func AddTask(name string, scanArea []string, companyID int64) (id int64, err error) {
	sql := `INSERT INTO task(name, scan_area,company_id) VALUES(?, ?, ?)`
	ret, err := db.Exec(sql, name, strings.Join(scanArea, ","), companyID)
	if err != nil {
		zap.L().Error("add task db.Exec error", zap.Error(err))
		return id, err
	}
	return ret.LastInsertId()
}

func GetTaskByID(id int64) (*models.Task, error) {
	sql := `SELECT * FROM task WHERE id = ?`
	var task models.Task
	if err := db.Get(&task, sql, id); err != nil {
		zap.L().Error("db.Get failed", zap.Error(err))
		return nil, err
	}
	return &task, nil
}

func GetTaskList(offset, count int) (task []*models.Task, err error) {
	sql := `SELECT task.*,company.name as company_name FROM task,company WHERE task.company_id=company.id LIMIT ?,?`
	if err := db.Select(&task, sql, offset, count); err != nil {
		zap.L().Error("db.Select failed", zap.Error(err))
		return nil, err
	}
	return task, nil
}

func DeleteTask(id int64) error {
	sql := `DELETE FROM task WHERE id = ?`
	_, err := db.Exec(sql, id)
	return err
}

func GetTaskCount(companyID int64) (count int64, err error) {
	sql := `SELECT count(id) from task where company_id = ?`
	err = db.Get(&count, sql, companyID)
	return count, err
}
