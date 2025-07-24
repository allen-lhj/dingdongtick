package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Label 标签模型
type Label struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`
	Name      string         `json:"name" gorm:"not null;size:100"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 不使用外键约束
	User  User   `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tasks []Task `json:"tasks,omitempty" gorm:"many2many:task_labels;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 指定表名
func (Label) TableName() string {
	return "labels"
}

// BeforeCreate GORM钩子，创建前生成UUID
func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
