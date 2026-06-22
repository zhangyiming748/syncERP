package csv

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"erp/model"
	"erp/storage"
)

//这里才应该是解析csv的代码
// 以D:\Users\Public\Github\syncERP\csv\test_output_utf8.csv文件为例执行解析 跳过空行（即全是逗号的行）
// 00114408 这一列为用户id  即数据表的EmployeeID列
// 王国庆 这一列为用户名  即数据表的EmployeeName列
// Z_COMMON_BIS 这一列为权限名  即数据表的EmployeeRole列
// 2025.05.24 这一列为创建时间  即数据表的CreatedAt列 需要解析为cst中国时间 然后插入
// 通用角色-接受内审管控 这一列为角色描述  即数据表EmployeeDescription列
// 解析出来的每一行组成一个EmployeeBefore20260413结构体 在数量允许的情况下批量插入数据库表

// ParseCSVAndImportToDB 解析 CSV 文件并导入到数据库（自动检测编码）
func ParseCSVAndImportToDB(filePath string) error {
	// 读取文件
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	var records [][]string

	// 尝试检测是否为 UTF-8 编码（简单检测：检查是否包含 BOM 或尝试直接解析）
	if isUTF8(fileData) {
		// 已经是 UTF-8，直接解析
		records, err = ParseCSVFromBytes(fileData)
		if err != nil {
			return fmt.Errorf("CSV解析失败: %w", err)
		}
	} else {
		// 需要转换编码
		utf8Data, err := ConvertGBKToUTF8(fileData)
		if err != nil {
			return fmt.Errorf("编码转换失败: %w", err)
		}
		records, err = ParseCSVFromBytes(utf8Data)
		if err != nil {
			return fmt.Errorf("CSV解析失败: %w", err)
		}
	}

	// 第三步：过滤空行并转换为结构体
	employees, err := convertRecordsToEmployees(records)
	if err != nil {
		return fmt.Errorf("数据转换失败: %w", err)
	}

	if len(employees) == 0 {
		return fmt.Errorf("没有有效数据")
	}

	fmt.Printf("共解析 %d 条有效记录\n", len(employees))

	// 第四步：批量插入数据库
	err = storage.BatchCreateEmployees(employees)
	if err != nil {
		return fmt.Errorf("批量插入数据库失败: %w", err)
	}

	fmt.Printf("成功插入 %d 条记录到数据库\n", len(employees))
	return nil
}

// convertRecordsToEmployees 将 CSV 记录转换为 EmployeeBefore20260413 结构体切片
func convertRecordsToEmployees(records [][]string) ([]model.EmployeeBefore20260413, error) {
	var employees []model.EmployeeBefore20260413

	for i, record := range records {
		// 跳过空行（所有字段都是空或只有逗号）
		if isEmptyRow(record) {
			continue
		}

		// 确保有足够的列
		if len(record) < 10 {
			fmt.Printf("警告: 第%d行字段数不足(%d个)，跳过\n", i+1, len(record))
			continue
		}

		// 解析创建时间 (格式: 2025.05.24)
		createdAt, err := parseChineseDate(record[7])
		if err != nil {
			fmt.Printf("警告: 第%d行时间解析失败(%s): %v，跳过\n", i+1, record[7], err)
			continue
		}

		// 构建结构体
		emp := model.EmployeeBefore20260413{
			EmployeeID:          strings.TrimSpace(record[1]),
			EmployeeName:        strings.TrimSpace(record[2]),
			EmployeeRole:        strings.TrimSpace(record[3]),
			EmployeeDescription: strings.TrimSpace(record[9]),
			CreatedAt:           createdAt,
			UpdatedAt:           time.Now(),
		}

		// 验证必要字段
		if emp.EmployeeID == "" {
			fmt.Printf("警告: 第%d行员工ID为空，跳过\n", i+1)
			continue
		}

		employees = append(employees, emp)
	}

	return employees, nil
}

// isEmptyRow 判断是否为空行（所有字段都是空或只包含空白字符）
func isEmptyRow(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

// parseChineseDate 解析中文日期格式 (2025.05.24) 为 CST 时间
func parseChineseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("日期字符串为空")
	}

	// 尝试解析格式: 2025.05.24
	t, err := time.ParseInLocation("2006.01.02", dateStr, time.FixedZone("CST", 8*3600))
	if err != nil {
		return time.Time{}, fmt.Errorf("日期格式错误: %w", err)
	}

	return t, nil
}

// isUTF8 检测字节数据是否为有效的 UTF-8 编码
func isUTF8(data []byte) bool {
	// 检查是否包含 UTF-8 BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return true
	}
	// 验证是否为有效的 UTF-8
	return utf8.Valid(data)
}
