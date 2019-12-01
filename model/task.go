package model

import "time"
import u "go-api/utils"

type Task struct {
	ID          uint `json:"id";gorm:"AUTO INCREMENT;PRIMARY KEY"`
	AccountID   uint `json:"user_id"`
	Account     Account
	Tasks       string    `json:"task";gorm:"size(150)"`
	Description string    `json:"description";gorm:"size(250)"`
	TaskDate    time.Time `json:"task_date"`
}

func (task *Task) AddTask() map[string]interface{} {

	if task.Tasks == "" {
		return u.Message(false, "missing task")
	}
	if task.Description == "" {
		return u.Message(false, "missing description")
	}
	//if task.TaskDate.IsZero(){
	//	return u.Message(false,"missing task date")
	//}

	GetDB().Create(task)
	if task.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	responds := u.Message(true, "Task Added Successfully!!")
	responds["data"] = task

	return responds

}
