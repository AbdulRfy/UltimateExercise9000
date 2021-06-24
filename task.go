package main

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	ID      uint32 `gorm:"size:100;not null;" json:"id"`
	UserId  uint32 `json:"userId,omitempty"`
	Name    string `json:"name"`
	DueDate string `json:"dueDate"`
}

type TaskAssign struct {
	gorm.Model

	Task         Task   `gorm:"constraint:OnDelete:CASCADE;"`
	TaskId       uint32 `gorm:"UNIQUE_INDEX:compositeindex;index;not null" json:"taskId"`
	AssigneEmail string `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"  json:"assigneEmail"`
}
