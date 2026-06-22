package main

import (
	"fmt"
	"log"
	"os"

	"erp/csv"
)

func main() {
	srcFile := "csv/EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"
	dstFile := "csv/EXPORT_20260413074548-昆仑集团含万方不含原生态权限_utf8.csv"

	fmt.Println("===========================================")
	fmt.Println("  ANSI CSV 转 UTF-8 工具")
	fmt.Println("===========================================")
	fmt.Println()

	// 检查源文件是否存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		log.Fatalf("错误：源文件不存在: %s", srcFile)
	}

	// 获取源文件信息
	srcInfo, _ := os.Stat(srcFile)
	fmt.Printf("源文件: %s\n", srcFile)
	fmt.Printf("大小:   %d bytes (%.2f MB)\n", srcInfo.Size(), float64(srcInfo.Size())/1024/1024)
	fmt.Println()

	// 执行转换
	fmt.Println("正在转换编码...")
	err := csv.ConvertANSICSVToUTF8File(srcFile, dstFile)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 获取目标文件信息
	dstInfo, _ := os.Stat(dstFile)
	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("✓ 转换成功！")
	fmt.Println("===========================================")
	fmt.Printf("目标文件: %s\n", dstFile)
	fmt.Printf("大小:     %d bytes (%.2f MB)\n", dstInfo.Size(), float64(dstInfo.Size())/1024/1024)
	fmt.Printf("增长:     %.2f%%\n", float64(dstInfo.Size()-srcInfo.Size())/float64(srcInfo.Size())*100)
	fmt.Println()
	fmt.Println("文件已保存，可以直接使用！")
}
