package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reminder 提醒模型
type Reminder struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TaskID    uuid.UUID      `json:"taskId" gorm:"type:uuid;not null;index"`
	RemindAt  time.Time      `json:"remindAt" gorm:"not null;index"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 不使用外键约束
	Task Task `json:"task,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName 指定表名
func (Reminder) TableName() string {
	return "reminders"
}

// BeforeCreate GORM钩子，创建前生成UUID
func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
