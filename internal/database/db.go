package database

import (
	"MediaTools/internal/config"

	"MediaTools/internal/models"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	switch config.DB.Type {
	case models.DBTypeSQLite:
		DB, err = gorm.Open(sqlite.Open(config.DB.DSN))

	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.DB.Type.String())
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	return AutoMigrate()
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.MediaTransferHistory{},
	)
}
