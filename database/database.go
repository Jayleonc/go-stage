package database

// 连接数据库，返回数据库连接

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func DBInstance() *gorm.DB {
	dsn := "root:jayleonc@tcp(127.0.0.1:3306)/go-stage?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {

		return nil
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB")
	return db
}
