package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ConvertGBKToUTF8 将 GBK/ANSI 编码的字节数据转换为 UTF-8
func ConvertGBKToUTF8(gbkData []byte) ([]byte, error) {
	reader := transform.NewReader(
		bytes.NewReader(gbkData),
		simplifiedchinese.GBK.NewDecoder(),
	)

	var utf8Buffer bytes.Buffer
	_, err := io.Copy(&utf8Buffer, reader)
	if err != nil {
		return nil, fmt.Errorf("GBK转UTF-8失败: %w", err)
	}

	return utf8Buffer.Bytes(), nil
}

// ConvertGBKFileToUTF8 将 GBK/ANSI 编码的文件转换为 UTF-8 字节数据
func ConvertGBKFileToUTF8(filePath string) ([]byte, error) {
	// 读取原始文件
	gbkData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 转换编码
	utf8Data, err := ConvertGBKToUTF8(gbkData)
	if err != nil {
		return nil, fmt.Errorf("编码转换失败: %w", err)
	}

	return utf8Data, nil
}

// ParseCSVFromBytes 从 UTF-8 字节数据解析 CSV
func ParseCSVFromBytes(utf8Data []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(utf8Data))

	// 配置 CSV 解析器选项
	reader.FieldsPerRecord = -1    // 允许每行字段数不同
	reader.LazyQuotes = true       // 允许非标准引号
	reader.TrimLeadingSpace = true // 去除前导空格

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}

	return records, nil
}

// ParseCSVFromFile 从文件路径解析 CSV（自动处理 GBK 到 UTF-8 转换）
func ParseCSVFromFile(filePath string) ([][]string, error) {
	// 第一步：将 GBK 文件转换为 UTF-8 字节数据
	utf8Data, err := ConvertGBKFileToUTF8(filePath)
	if err != nil {
		return nil, fmt.Errorf("编码转换失败: %w", err)
	}

	// 第二步：解析 UTF-8 CSV 数据
	records, err := ParseCSVFromBytes(utf8Data)
	if err != nil {
		return nil, fmt.Errorf("CSV解析失败: %w", err)
	}

	return records, nil
}

// ParseCSVWithHeader 解析带表头的 CSV，返回表头和记录
func ParseCSVWithHeader(filePath string) (header []string, records [][]string, err error) {
	allRecords, err := ParseCSVFromFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	if len(allRecords) == 0 {
		return nil, nil, fmt.Errorf("CSV文件为空")
	}

	// 第一行作为表头
	header = allRecords[0]
	// 剩余行作为数据记录
	records = allRecords[1:]

	return header, records, nil
}
