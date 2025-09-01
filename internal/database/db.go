package database

import (
	"MediaTools/internal/config"
	"strings"

	"MediaTools/internal/models"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	var err error
	switch strings.ToLower(config.DB.Type) {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DB.DSN))

	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.DB.Type)
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	return AutoMigrate()
}

func AutoMigrate() error {
	return db.AutoMigrate(
		&models.MediaTransferHistory{},
	)
}
