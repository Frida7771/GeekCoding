package models

import "gorm.io/gorm"

type Submit_Basic struct {
	gorm.Model
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	UserIdentity    string `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	Path            string `gorm:"column:path;type:varchar(255)" json:"path"`
}

func (table *Submit_Basic) TableName() string {
	return "submit_basic"
}
