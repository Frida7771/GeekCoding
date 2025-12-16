package models

import "gorm.io/gorm"

type ProblemBasic struct {
	gorm.Model
	Identity         string             `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemCategorys []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	CategoryID       uint               `gorm:"column:category_id;type:varchar(255);" json:"category_id"`
	Title            string             `gorm:"column:title;type:varchar(255);" json:"title"`
	Content          string             `gorm:"column:content;type:text;" json:"content"`
	MaxRuntime       int                `gorm:"column:max_runtime;type:int(11)" json:"max_runtime"`
	MaxMem           int                `gorm:"column:max_mem;type:int(11)" json:"max_mem"`
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList_Basic(keyword string, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategorys").Preload("ProblemCategorys.CategoryBasic").
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("Right Join problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id =(SELECT id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}
