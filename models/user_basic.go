package models

import "gorm.io/gorm"

type User_Basic struct {
	gorm.Model
	Identity         string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name             string `gorm:"column:username;type:varchar(100)" json:"username"`
	Password         string `gorm:"column:password;type:varchar(32)" json:"password"`
	Phone            string `gorm:"column:phone;type:varchar(20)" json:"phone"`
	Email            string `gorm:"column:email;type:varchar(100)" json:"email"`
	FinishProblemNum int    `gorm:"column:finish_problem_num;type:int(11)" json:"finish_problem_num"`
	SubmitNum        int    `gorm:"column:submit_num;type:int(11)" json:"submit_num"`
	IsAdmin          int    `gorm:"column:is_admin;type:tinyint(1)" json:"is_admin"`
}

func (table *User_Basic) TableName() string {
	return "user_basic"
}
