package csv

import (
	"fmt"
	"os"
)

// ConvertANSICSVToUTF8File 将 ANSI 编码的 CSV 文件转换为 UTF-8 编码并保存为新文件
// srcPath: 源文件路径（ANSI/GBK 编码）
// dstPath: 目标文件路径（UTF-8 编码）
func ConvertANSICSVToUTF8File(srcPath, dstPath string) error {
	// 第一步：读取 ANSI 文件并转换为 UTF-8 字节数据
	utf8Data, err := ConvertGBKFileToUTF8(srcPath)
	if err != nil {
		return fmt.Errorf("编码转换失败: %w", err)
	}

	// 第二步：将 UTF-8 数据写入新文件
	err = os.WriteFile(dstPath, utf8Data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
