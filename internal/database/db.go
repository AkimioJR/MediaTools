package database

import (
	"MediaTools/internal/config"
	"strings"

	"MediaTools/internal/models"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
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
	logrus.Debug("数据库连接成功，正在进行自动迁移...")
	return AutoMigrate()
}

func AutoMigrate() error {
	return db.AutoMigrate(
		&models.MediaTransferHistory{},
	)
}
