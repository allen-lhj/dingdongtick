package dal

import (
	"fmt"
	"log"

	"ticktick-backend/config"
	"ticktick-backend/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接管理器
type Database struct {
	GORM *gorm.DB
	SQLX *sqlx.DB
}

// NewDatabase 创建新的数据库连接
func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.GetDSN()

	// 初始化GORM连接 - 禁用外键约束
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		return nil, fmt.Errorf("连接GORM数据库失败: %w", err)
	}

	// 初始化SQLX连接
	sqlxDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接SQLX数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	db := &Database{
		GORM: gormDB,
		SQLX: sqlxDB,
	}

	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func (db *Database) AutoMigrate() error {
	log.Println("开始数据库迁移...")

	// 按依赖顺序迁移表
	err := db.GORM.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Label{},
		&models.Task{},
		&models.Reminder{},
		&models.TaskRecurrenceException{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 创建索引
	if err := db.createIndexes(); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	log.Println("数据库迁移完成")
	return nil
}

// createIndexes 创建必要的索引
func (db *Database) createIndexes() error {
	// 用户邮箱唯一索引
	if err := db.GORM.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 项目名称在用户内唯一
	if err := db.GORM.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_user_name ON projects(user_id, name) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 标签名称在用户内唯一
	if err := db.GORM.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_labels_user_name ON labels(user_id, name) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 循环任务例外的唯一索引
	if err := db.GORM.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_task_exceptions_unique ON task_recurrence_exceptions(recurring_task_id, original_time) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 性能索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_due_time ON tasks(due_time)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority)",
		"CREATE INDEX IF NOT EXISTS idx_reminders_remind_at ON reminders(remind_at)",
	}

	for _, indexSQL := range indexes {
		if err := db.GORM.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

// Close 关闭数据库连接
func (db *Database) Close() error {
	if db.SQLX != nil {
		if err := db.SQLX.Close(); err != nil {
			return err
		}
	}

	if db.GORM != nil {
		sqlDB, err := db.GORM.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}

// Ping 测试数据库连接
func (db *Database) Ping() error {
	if db.SQLX != nil {
		return db.SQLX.Ping()
	}
	return nil
}
