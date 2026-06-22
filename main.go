package main

import (
	"fmt"
	"log"

	"erp/csv"
	"erp/storage"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("  CSV 数据导入工具")
	fmt.Println("===========================================")
	fmt.Println()

	// 初始化数据库
	dbPath := "C:\\Users\\zen\\Documents\\erp.db"
	fmt.Printf("初始化数据库: %s\n", dbPath)
	storage.SetSqlite(dbPath)

	// 同步表结构
	fmt.Println("同步表结构...")
	err := storage.SyncEmployeeTable()
	if err != nil {
		log.Fatalf("表结构同步失败: %v", err)
	}
	fmt.Println("✓ 表结构同步完成")
	fmt.Println()

	// 解析并导入 CSV
	csvFile := "csv/test_output_utf8.csv"
	fmt.Printf("开始解析 CSV 文件: %s\n", csvFile)
	fmt.Println()

	err = csv.ParseCSVAndImportToDB(csvFile)
	if err != nil {
		log.Fatalf("导入失败: %v", err)
	}

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("✓ 所有操作完成！")
	fmt.Println("===========================================")
}
