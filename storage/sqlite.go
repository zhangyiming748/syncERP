package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var gormDB *gorm.DB

func SetSqlite(path string) *gorm.DB {
	// 创建数据目录
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		log.Fatal("无法创建数据目录:", err)
	}

	// 使用纯Go SQLite驱动连接数据库
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
			NoLowerCase:   true, // 不转换为小写，保持驼峰式命名
		},
	})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	gormDB = db

	log.Println("成功连接到SQLite数据库")
	return gormDB
}

func GetSqlite() *gorm.DB {
	return gormDB
}
