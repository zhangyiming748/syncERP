package model

import (
	"time"

	"gorm.io/gorm"
)

type EmployeeBefore20260413 struct {
	EmployeeID          string         `gorm:"column:employee_id;type:text" json:"employee_id"`
	EmployeeName        string         `gorm:"column:employee_name;type:text" json:"employee_name"`
	EmployeeRole        string         `gorm:"column:employee_role;type:text" json:"employee_role"`
	EmployeeDescription string         `gorm:"column:employee_description;type:text" json:"employee_description"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
