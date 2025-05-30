package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

var MysqlDB *gorm.DB

func InitMySQL() error {
	slog.Info("初始化mysql")
	// 获取反序列化后的mysql
	mysqlConfig := AppConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DBName, mysqlConfig.Charset)

	// 连接mysql
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出所有 SQL
	})
	if err != nil {
		return fmt.Errorf("mysql连接失败: %w", err)
	}
	slog.Info("mysql连接成功")
	return nil
}
