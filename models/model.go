package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Avatar   string
	Projects []Project   `gorm:"many2many:user_projects;"` // Projects associated with this user
	UserTask []TaskModel `gorm:"foreignKey:UserTaskID"`
}
type Project struct {
	gorm.Model
	ProjectName string      `json:"project_name" binding:"required"`
	Users       []User      `gorm:"many2many:user_projects;"` // Users associated with this project
	Task        []TaskModel `gorm:"foreignKey:ProjectTaskID"`
}

type TaskModel struct {
	gorm.Model
	TaskName      string         `json:"task_name"`
	Status        string         `json:"status" gorm:"default:in-queue"`
	Description   string         `json:"description"`
	ProjectTaskID uint           `json:"project_id"`
	SubTask       []SubTaskModel `gorm:"foreignKey:TaskSubTaskID"`
	UserTaskID    uint           `json:"user_task_id" gorm:"default:null"`
}

type SubTaskModel struct {
	gorm.Model
	TaskName      string `json:"task_name"`
	Status        string `json:"status" gorm:"default:in-queue"`
	Description   string `json:"description"`
	TaskSubTaskID uint   `json:"task_sub_task_id"`
}
