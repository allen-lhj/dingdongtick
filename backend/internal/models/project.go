package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`
	Name      string         `json:"name" gorm:"not null;size:255"`
	Color     string         `json:"color" gorm:"not null;size:7;default:#CCCCCC"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 不使用外键约束
	User  User   `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tasks []Task `json:"tasks,omitempty" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate GORM钩子，创建前生成UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
