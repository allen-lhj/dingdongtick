package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskRecurrenceException 循环任务例外模型
type TaskRecurrenceException struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RecurringTaskID uuid.UUID      `json:"recurringTaskId" gorm:"type:uuid;not null;index"`
	OriginalTime    time.Time      `json:"originalTime" gorm:"not null"`         // 原始实例的时间
	NewTaskID       *uuid.UUID     `json:"newTaskId,omitempty" gorm:"type:uuid"` // 如果被修改，指向新任务
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 不使用外键约束
	RecurringTask Task  `json:"recurringTask,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NewTask       *Task `json:"newTask,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName 指定表名
func (TaskRecurrenceException) TableName() string {
	return "task_recurrence_exceptions"
}

// BeforeCreate GORM钩子，创建前生成UUID
func (tre *TaskRecurrenceException) BeforeCreate(tx *gorm.DB) error {
	if tre.ID == uuid.Nil {
		tre.ID = uuid.New()
	}
	return nil
}

// IsDeleted 判断是否为删除的例外（没有新任务ID）
func (tre *TaskRecurrenceException) IsDeleted() bool {
	return tre.NewTaskID == nil
}

// IsModified 判断是否为修改的例外（有新任务ID）
func (tre *TaskRecurrenceException) IsModified() bool {
	return tre.NewTaskID != nil
}
