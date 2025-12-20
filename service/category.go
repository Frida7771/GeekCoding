package service

import (
	"GeekCoding/define"
	"GeekCoding/help"
	"GeekCoding/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategoryList
// @Tags         Admin Method
// @Summary      Get Category List
// @Param        Authorization header     string     true  "Authorization"
// @Param        page  query     int     false  "page"
// @Param        size  query     int     false  "size"
// @Param        keyword  query     string     false  "keyword"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /category-list [get]
func GetCategoryList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetCategoryList Page strconv Error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	keyword := c.Query("keyword")

	categoryList := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(new(models.CategoryBasic)).Where("name LIKE ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&categoryList).Error
	if err != nil {
		log.Println("GetCategoryList Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get category list error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  categoryList,
			"count": count,
		},
	})
}

// CreateCategory
// @Tags         Admin Method
// @Summary      Create Category
// @Param        Authorization header     string     true  "Authorization"
// @Param        name  formData     string     true  "name"
// @Param        parent_id  formData     int     false  "parent_id"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /category-create [post]
func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	parent_id, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := &models.CategoryBasic{
		Identity: help.GetUUID(),
		Name:     name,
		ParentID: uint(parent_id),
	}
	err := models.DB.Create(category).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create category error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "create category success",
	})

}

// DeleteCategory
// @Tags         Admin Method
// @Summary      Delete Category
// @Param        Authorization header     string     true  "Authorization"
// @Param        identity  query     string     true  "identity"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /category-delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "identity is required",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id =(SELECT id FROM category_basic WHERE identity = ? )", identity).Count(&cnt).Error
	if err != nil {
		log.Println("Get Problem Category Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get problem category error",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "category has problem, delete failed",
		})
		return
	}
	err = models.DB.Unscoped().Where("identity = ?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		log.Println("Delete Category Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "delete category error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "delete category success",
	})
}

// UpdateCategory
// @Tags         Admin Method
// @Summary      Update Category
// @Param        Authorization header     string     true  "Authorization"
// @Param        parent_id  formData     int     false  "parent_id"
// @Param        identity  query     string     true  "identity"
// @Param        name  formData     string     true  "name"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /category-update [put]
func UpdateCategory(c *gin.Context) {
	identity := c.Query("identity")
	parent_id, _ := strconv.Atoi(c.PostForm("parent_id"))
	name := c.PostForm("name")
	if name == "" || identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "name and identity are required",
		})
		return
	}

	category := &models.CategoryBasic{
		Identity: identity,
		ParentID: uint(parent_id),
		Name:     name,
	}
	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "update category error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "update category success",
	})

}
