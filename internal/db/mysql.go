package db

import (
	"fmt"

	"github.com/SideSwapIN/Analystic/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

// InitDB 初始化数据库连接
func InitMysqlDB() error {
	db, err := gorm.Open(mysql.Open(config.GetConfig().Database.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect mysql database error: %v", err)
	}

	// Set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to set mysql connection pool error: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	MysqlDB = db
	return nil
}
