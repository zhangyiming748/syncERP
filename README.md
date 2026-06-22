# syncERP - CSV 数据导入工具

## 📋 项目简介

syncERP 是一个用于处理企业 ERP 系统员工权限数据的工具，支持自动检测 CSV 文件编码（UTF-8/GBK/GB2312/GB18030），解析数据并批量导入 SQLite 数据库。

## ✨ 核心特性

- 🔍 **智能编码检测**：自动识别 UTF-8、GBK、GB2312、GB18030 等编码格式
- 🔄 **自动编码转换**：非 UTF-8 编码自动转换为 UTF-8，避免中文乱码
- 📊 **高效批量导入**：每批 100 条记录批量插入，提升性能
- 🛡️ **数据验证**：自动跳过空行、验证必要字段、处理异常数据
- 🕐 **时区支持**：正确解析中国标准时间（CST, UTC+8）

## 🏗️ 项目结构

```
syncERP/
├── csv/                          # CSV 处理模块
│   ├── parse.go                  # CSV 解析与数据库导入核心逻辑
│   ├── ansi2utf8.go              # ANSI/GBK 到 UTF-8 编码转换
│   └── ansi2utf8_test.go         # 编码转换单元测试
├── model/                        # 数据模型
│   └── employee_before_20260413.go  # 员工数据结构定义
├── storage/                      # 数据存储层
│   ├── sqlite.go                 # SQLite 数据库连接管理
│   └── employee_before_20260413.go  # 员工数据 CRUD 操作
├── main.go                       # 主程序入口
├── go.mod                        # Go 模块依赖配置
└── README.md                     # 项目文档
```

## 🚀 快速开始

### 前置要求

- Go 1.20+
- Windows/Linux/macOS

### 安装依赖

```bash
go mod tidy
```

### 运行程序

```bash
# 编译
go build -o import_csv.exe .

# 运行
./import_csv.exe
```

## 📖 工作流程

### 完整流程图

```
┌─────────────────────┐
│   启动程序 (main)    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  初始化 SQLite 数据库 │
│  (data/erp.db)       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  同步表结构           │
│  (AutoMigrate)       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  读取 CSV 文件        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  检测文件编码         │
│  isUTF8()            │
└────┬────────────┬───┘
     │            │
     │ UTF-8      │ 非 UTF-8
     │            │ (GBK/GB2312/GB18030)
     ▼            ▼
┌──────────┐  ┌──────────────────┐
│ 直接解析  │  │ 转换为 UTF-8      │
│ CSV      │  │ ConvertGBKToUTF8 │
└────┬─────┘  └────────┬─────────┘
     │                 │
     └────────┬────────┘
              │
              ▼
┌─────────────────────┐
│  过滤空行            │
│  isEmptyRow()        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  解析字段映射         │
│  • EmployeeID       │
│  • EmployeeName     │
│  • EmployeeRole     │
│  • CreatedAt (CST)  │
│  • EmployeeDesc     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  批量插入数据库       │
│  (每批 100 条)       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  完成                │
└─────────────────────┘
```

### 详细步骤说明

#### 1. 编码检测 (`isUTF8`)

```go
// 检测逻辑：
// 1. 检查 UTF-8 BOM (EF BB BF)
// 2. 使用 utf8.Valid() 验证是否为有效 UTF-8 编码
func isUTF8(data []byte) bool {
    if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
        return true  // 有 BOM 标记
    }
    return utf8.Valid(data)  // 验证 UTF-8 有效性
}
```

**检测结果：**
- ✅ **UTF-8**：直接解析，无需转换
- ⚠️ **非 UTF-8**：判定为 GBK/GB2312/GB18030，需要转换

#### 2. 编码转换 (`ConvertGBKToUTF8`)

```go
// 使用 golang.org/x/text 进行实时流式转换
func ConvertGBKToUTF8(gbkData []byte) ([]byte, error) {
    reader := transform.NewReader(
        bytes.NewReader(gbkData),
        simplifiedchinese.GBK.NewDecoder(),
    )
    
    var utf8Buffer bytes.Buffer
    io.Copy(&utf8Buffer, reader)
    
    return utf8Buffer.Bytes(), nil
}
```

**支持的编码：**
- GBK
- GB2312
- GB18030

#### 3. CSV 解析 (`ParseCSVFromBytes`)

```go
// 配置选项
reader.FieldsPerRecord = -1      // 允许每行字段数不同
reader.LazyQuotes = true         // 允许非标准引号
reader.TrimLeadingSpace = true   // 去除前导空格
```

#### 4. 数据过滤与转换 (`convertRecordsToEmployees`)

**空行检测：**
```go
func isEmptyRow(record []string) bool {
    for _, field := range record {
        if strings.TrimSpace(field) != "" {
            return false  // 有非空字段，不是空行
        }
    }
    return true  // 所有字段都为空
}
```

**字段映射：**

| CSV 列索引 | 字段内容 | 结构体字段 | 说明 |
|-----------|---------|-----------|------|
| 1 | 00114408 | EmployeeID | 用户 ID（必填） |
| 2 | 王国庆 | EmployeeName | 用户名 |
| 3 | Z_COMMON_BIS | EmployeeRole | 权限名 |
| 7 | 2025.05.24 | CreatedAt | 创建时间（CST） |
| 9 | 通用角色-接受内审管控 | EmployeeDescription | 角色描述 |

**时间解析：**
```go
// 格式：2025.05.24 → CST (UTC+8)
func parseChineseDate(dateStr string) (time.Time, error) {
    return time.ParseInLocation("2006.01.02", dateStr, 
        time.FixedZone("CST", 8*3600))
}
```

#### 5. 批量插入 (`BatchCreateEmployees`)

```go
// GORM 批量插入，每批 100 条
db.CreateInBatches(employees, 100)
```

**优势：**
- 减少数据库交互次数
- 提升插入性能（相比逐条插入快 10-50 倍）
- 自动事务管理

## 📊 数据模型

### EmployeeBefore20260413

```go
type EmployeeBefore20260413 struct {
    EmployeeID          string         `gorm:"column:employee_id;type:text"`
    EmployeeName        string         `gorm:"column:employee_name;type:text"`
    EmployeeRole        string         `gorm:"column:employee_role;type:text"`
    EmployeeDescription string         `gorm:"column:employee_description;type:text"`
    CreatedAt           time.Time      // 创建时间（CST）
    UpdatedAt           time.Time      // 更新时间
    DeletedAt           gorm.DeletedAt // 软删除时间戳
}
```

**数据库配置：**
- 单数表名：`employee_before_20260413`
- 驼峰命名：保持原样，不转换为小写
- 软删除：使用 GORM DeletedAt

## 🧪 测试

### 运行单元测试

```bash
cd csv
go test -v
```

### 测试覆盖

- ✅ ANSI 到 UTF-8 文件转换
- ✅ GBK 字节数据转换
- ✅ CSV 解析功能
- ✅ 带表头 CSV 解析
- ✅ 内容完整性验证

### 测试输出示例

```
=== RUN   TestConvertANSICSVToUTF8File
    ansi2utf8_test.go:42: ✓ 转换成功
    ansi2utf8_test.go:43:   源文件大小: 11012892 bytes
    ansi2utf8_test.go:44:   目标文件大小: 13603706 bytes
--- PASS: TestConvertANSICSVToUTF8File (0.05s)
```

## 📝 使用示例

### 基本用法

```go
package main

import (
    "log"
    "erp/csv"
    "erp/storage"
)

func main() {
    // 1. 初始化数据库
    storage.SetSqlite("data/erp.db")
    
    // 2. 同步表结构
    storage.SyncEmployeeTable()
    
    // 3. 解析并导入 CSV（自动检测编码）
    err := csv.ParseCSVAndImportToDB("data.csv")
    if err != nil {
        log.Fatal(err)
    }
}
```

### 单独使用编码转换

```go
// 方法 1：转换文件并保存
err := csv.ConvertANSICSVToUTF8File(
    "input_ansi.csv",
    "output_utf8.csv",
)

// 方法 2：转换后直接解析
records, err := csv.ParseCSVFromFile("input_ansi.csv")

// 方法 3：分步处理
utf8Data, _ := csv.ConvertGBKFileToUTF8("input.csv")
records, _ := csv.ParseCSVFromBytes(utf8Data)
```

## ⚙️ 配置说明

### 数据库配置

```go
// storage/sqlite.go
db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        SingularTable: true, // 使用单数表名
        NoLowerCase:   true, // 保持驼峰命名
    },
})
```

### CSV 解析配置

```go
// csv/parse.go
reader.FieldsPerRecord = -1      // 灵活字段数
reader.LazyQuotes = true         // 宽松引号处理
reader.TrimLeadingSpace = true   // 自动去空格
```

## 🔧 常见问题

### Q1: 为什么需要编码转换？

**A:** Go 标准库 `encoding/csv` 默认期望 UTF-8 编码。ANSI/GBK 编码的中文字符会被错误解析为乱码。

**示例：**
```
原始数据：张三
未转换：å¼ ä¸‰ (乱码)
转换后：张三 (正确)
```

### Q2: UTF-8 文件比 GBK 大多少？

**A:** 大约 20-30%。原因：
- GBK：中文字符占 2 字节
- UTF-8：中文字符占 3 字节

**实测数据：**
- GBK: 11,012,892 bytes (10.5 MB)
- UTF-8: 13,603,706 bytes (13.0 MB)
- 增长：23.5%

### Q3: 如何处理超大 CSV 文件？

**A:** 当前实现已使用批量插入（每批 100 条），对于 90,000+ 记录的文件，导入时间在几秒内完成。如需处理更大文件，可以：
1. 增加批次大小（如 500 或 1000）
2. 使用流式解析，避免一次性加载全部数据到内存

### Q4: 如何验证导入的数据？

```go
// 查询所有记录
employees, err := storage.ListEmployees()

// 根据 ID 查询
emp, err := storage.GetEmployeeByID("00114408")

// 统计数量
var count int64
db.Model(&model.EmployeeBefore20260413{}).Count(&count)
```

## 📦 依赖项

```go
require (
    github.com/glebarez/sqlite v1.5.0
    golang.org/x/text v0.38.0
    gorm.io/gorm v1.25.0
)
```

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📮 联系方式

如有问题或建议，请提交 Issue。
