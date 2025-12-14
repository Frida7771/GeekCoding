package test

import (
	"fmt"
	"testing"

	"GeekCoding/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGorm(t *testing.T) {
	dsn := "root:#Etnlhy1396917302@tcp(127.0.0.1:3306)/Geek_Coding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	data := make([]*models.Problem, 0)
	err = db.Find(&data).Error
	if err != nil {
		panic(err)
	}

	for _, v := range data {
		fmt.Printf("Problem: %+v\n", v)

	}

}
