package storage

import (
	"erp/model"
	"fmt"
)

// SyncEmployeeTable 同步员工表结构
func SyncEmployeeTable() error {
	db := GetSqlite()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	err := db.AutoMigrate(&model.EmployeeBefore20260413{})
	if err != nil {
		return err
	}
	return nil
}

// CreateEmployee 创建员工记录
func CreateEmployee(emp *model.EmployeeBefore20260413) error {
	db := GetSqlite()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.Create(emp).Error
}

// GetEmployeeByID 根据员工ID查询员工
func GetEmployeeByID(employeeID string) (*model.EmployeeBefore20260413, error) {
	db := GetSqlite()
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	var emp model.EmployeeBefore20260413
	err := db.Where("employee_id = ?", employeeID).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// ListEmployees 查询所有员工列表
func ListEmployees() ([]model.EmployeeBefore20260413, error) {
	db := GetSqlite()
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	var employees []model.EmployeeBefore20260413
	err := db.Find(&employees).Error
	return employees, err
}

// UpdateEmployee 更新员工信息
func UpdateEmployee(emp *model.EmployeeBefore20260413) error {
	db := GetSqlite()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.Save(emp).Error
}

// DeleteEmployee 删除员工记录(软删除)
func DeleteEmployee(employeeID string) error {
	db := GetSqlite()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.Where("employee_id = ?", employeeID).Delete(&model.EmployeeBefore20260413{}).Error
}

// BatchCreateEmployees 批量创建员工记录
func BatchCreateEmployees(employees []model.EmployeeBefore20260413) error {
	db := GetSqlite()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.CreateInBatches(employees, 100).Error
}
