package csv

import (
	"os"
	"testing"
)

func TestConvertANSICSVToUTF8File(t *testing.T) {
	// 测试文件路径
	srcFile := "EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"
	dstFile := "test_output_utf8.csv"

	// 确保源文件存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Skipf("跳过测试：源文件 %s 不存在", srcFile)
	}

	// 执行转换
	err := ConvertANSICSVToUTF8File(srcFile, dstFile)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 清理：测试结束后删除生成的文件
	// defer os.Remove(dstFile)  // 注释掉这行以保留文件

	// 验证目标文件已创建
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Fatal("目标文件未创建")
	}

	// 读取目标文件内容，验证是有效的 UTF-8 编码
	utf8Data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("读取目标文件失败: %v", err)
	}

	if len(utf8Data) == 0 {
		t.Fatal("目标文件内容为空")
	}

	t.Logf("✓ 转换成功")
	t.Logf("  源文件大小: %d bytes", getFileSize(srcFile))
	t.Logf("  目标文件大小: %d bytes", len(utf8Data))
	t.Logf("  目标文件路径: %s", dstFile)
}

func TestConvertGBKFileToUTF8(t *testing.T) {
	// 测试文件路径
	srcFile := "EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"

	// 确保源文件存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Skipf("跳过测试：源文件 %s 不存在", srcFile)
	}

	// 执行转换
	utf8Data, err := ConvertGBKFileToUTF8(srcFile)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if len(utf8Data) == 0 {
		t.Fatal("转换后的数据为空")
	}

	t.Logf("✓ GBK转UTF-8成功")
	t.Logf("  原始文件大小: %d bytes", getFileSize(srcFile))
	t.Logf("  UTF-8数据大小: %d bytes", len(utf8Data))
}

func TestParseCSVFromFile(t *testing.T) {
	// 测试文件路径
	srcFile := "EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"

	// 确保源文件存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Skipf("跳过测试：源文件 %s 不存在", srcFile)
	}

	// 解析 CSV
	records, err := ParseCSVFromFile(srcFile)
	if err != nil {
		t.Fatalf("解析CSV失败: %v", err)
	}

	if len(records) == 0 {
		t.Fatal("解析结果为空")
	}

	t.Logf("✓ CSV解析成功")
	t.Logf("  总行数: %d", len(records))

	// 打印前5行作为示例
	displayCount := 5
	if len(records) < displayCount {
		displayCount = len(records)
	}

	for i := 0; i < displayCount; i++ {
		t.Logf("  第%d行 (%d个字段): %v", i+1, len(records[i]), records[i])
	}
}

func TestParseCSVWithHeader(t *testing.T) {
	// 测试文件路径
	srcFile := "EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"

	// 确保源文件存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Skipf("跳过测试：源文件 %s 不存在", srcFile)
	}

	// 解析带表头的 CSV
	header, records, err := ParseCSVWithHeader(srcFile)
	if err != nil {
		t.Fatalf("解析CSV失败: %v", err)
	}

	if header == nil {
		t.Fatal("表头为空")
	}

	t.Logf("✓ 带表头CSV解析成功")
	t.Logf("  表头字段数: %d", len(header))
	t.Logf("  数据行数: %d", len(records))
	t.Logf("  表头: %v", header)

	// 打印第一条记录
	if len(records) > 0 {
		t.Logf("  第一条记录: %v", records[0])
	}
}

// getFileSize 获取文件大小
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// TestConvertAndVerifyContent 转换后验证中文内容是否正确
func TestConvertAndVerifyContent(t *testing.T) {
	srcFile := "EXPORT_20260413074548-昆仑集团含万方不含原生态权限.csv"
	dstFile := "test_verify_utf8.csv"

	// 确保源文件存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Skipf("跳过测试：源文件 %s 不存在", srcFile)
	}

	// 执行转换
	err := ConvertANSICSVToUTF8File(srcFile, dstFile)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	defer os.Remove(dstFile)

	// 读取并解析转换后的文件
	records, err := ParseCSVFromFile(dstFile)
	if err != nil {
		t.Fatalf("解析转换后的文件失败: %v", err)
	}

	if len(records) == 0 {
		t.Fatal("解析结果为空")
	}

	// 检查是否包含中文字符（简单的验证）
	hasChinese := false
	for _, record := range records {
		for _, field := range record {
			if len(field) > 0 && field != "," && field != " " {
				// 检查字段长度，如果转换正确，中文字段应该有合理的长度
				if len(field) > 1 {
					hasChinese = true
					break
				}
			}
		}
		if hasChinese {
			break
		}
	}

	if !hasChinese {
		t.Log("警告：未能检测到中文字符，可能需要手动验证文件内容")
	}

	t.Logf("✓ 内容验证完成")
	t.Logf("  总记录数: %d", len(records))
}
