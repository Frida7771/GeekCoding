package models

import "gorm.io/gorm"

type Submit_Basic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`
	UserBasic       *User_Basic   `gorm:"foreignKey:identity;references:user_identity"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	Path            string        `gorm:"column:path;type:varchar(255)" json:"path"`
	Status          int           `gorm:"column:status;type:tinyint(1)" json:"status"`
}

func (table *Submit_Basic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(Submit_Basic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")

	if problemIdentity != "" {
		tx = tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx = tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}

	return tx
}
