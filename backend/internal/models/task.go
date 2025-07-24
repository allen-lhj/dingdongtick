package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusIncomplete TaskStatus = "incomplete"
	TaskStatusCompleted  TaskStatus = "completed"
)

// Task 任务模型
type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`
	ProjectID   uuid.UUID      `json:"projectId" gorm:"type:uuid;not null;index"`
	ParentID    *uuid.UUID     `json:"parentId,omitempty" gorm:"type:uuid;index"` // 子任务的父任务ID
	Title       string         `json:"title" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TaskStatus     `json:"status" gorm:"not null;default:incomplete;check:status IN ('incomplete', 'completed')"`
	Priority    int            `json:"priority" gorm:"not null;default:4;check:priority BETWEEN 1 AND 4"`
	StartTime   *time.Time     `json:"startTime,omitempty"`
	DueTime     *time.Time     `json:"dueTime,omitempty"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
	RRuleString string         `json:"rruleString,omitempty" gorm:"type:text"` // RFC 5545 循环规则
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 不使用外键约束
	User       User                      `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Project    Project                   `json:"project,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent     *Task                     `json:"parent,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SubTasks   []Task                    `json:"subTasks,omitempty" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Labels     []Label                   `json:"labels,omitempty" gorm:"many2many:task_labels;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reminders  []Reminder                `json:"reminders,omitempty" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Exceptions []TaskRecurrenceException `json:"exceptions,omitempty" gorm:"foreignKey:RecurringTaskID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate GORM钩子，创建前生成UUID
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsRecurring 判断是否为循环任务
func (t *Task) IsRecurring() bool {
	return t.RRuleString != ""
}

// IsCompleted 判断任务是否已完成
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// MarkCompleted 标记任务为已完成
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
}

// MarkIncomplete 标记任务为未完成
func (t *Task) MarkIncomplete() {
	t.Status = TaskStatusIncomplete
	t.CompletedAt = nil
}
