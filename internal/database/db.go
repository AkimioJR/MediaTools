package database

import (
	"MediaTools/internal/config"
	"MediaTools/internal/model"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	switch config.DB.Type {
	case model.DBTypeSQLite:
		DB, err = gorm.Open(sqlite.Open("gorm.db"))

	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.DB.Type.String())
	}

	return err
}
