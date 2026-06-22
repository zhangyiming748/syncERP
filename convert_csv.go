package main

import (
	"fmt"
	"log"

	"erp/csv"
)

func main() {
	srcFile := "csv/EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"
	dstFile := "csv/EXPORT_20260413074548-昆仑集团含万方不含原生态权限_utf8.csv"

	fmt.Println("开始转换 ANSI 编码的 CSV 文件为 UTF-8...")

	err := csv.ConvertANSICSVToUTF8File(srcFile, dstFile)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	fmt.Printf("✓ 转换成功！\n")
	fmt.Printf("  源文件: %s (ANSI/GBK)\n", srcFile)
	fmt.Printf("  目标文件: %s (UTF-8)\n", dstFile)
}
