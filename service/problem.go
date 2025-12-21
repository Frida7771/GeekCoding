package service

import (
	"GeekCoding/define"
	"GeekCoding/help"
	"GeekCoding/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProblemList
// @Tags         Public Method
// @Summary      Get Problem List
// @Param        page  query     int     false  "page"
// @Param        size  query     int     false  "size"
// @Param        keyword  query     string     false  "keyword"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	list := make([]*models.ProblemBasic, 0)
	err = models.GetProblemList_Basic(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error
	if err != nil {
		log.Println("GetProblemList Count Error:", err)
		return
	}
	err = models.GetProblemList_Basic(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Tags         Public Method
// @Summary      Get Problem Detail
// @Param        identity  query     string     true  "problem identity"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "problem identity is required",
		})
		return
	}

	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategorys").Preload("ProblemCategorys.CategoryBasic").First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "problem not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "get problem detail error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// ProblemCreate
// @Tags         Admin Method
// @Summary      Problem Create
// @Param        Authorization header     string     true  "Authorization"
// @Param        title  formData     string     true  "title"
// @Param        content  formData     string     true  "content"
// @Param        max_runtime  formData     int       true  "max_runtime"
// @Param        max_mem  formData     int       true  "max_mem"
// @Param        category_ids  formData     string     true  "category_ids"
// @Param        test_cases  body     string     true  "test_cases"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	max_runtime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	max_mem, _ := strconv.Atoi(c.PostForm("max_mem"))
	category_ids := c.PostFormArray("category_ids")
	test_cases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(category_ids) == 0 || len(test_cases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "parameters are required",
		})
		return
	}
	identity := help.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: max_runtime,
		MaxMem:     max_mem,
	}

	//deal category
	categoryBasic := make([]*models.ProblemCategory, 0)
	for _, id := range category_ids {
		category_id, _ := strconv.Atoi(id)
		categoryBasic = append(categoryBasic, &models.ProblemCategory{
			CategoryID: uint(category_id),
		})
	}
	data.ProblemCategorys = categoryBasic

	//deal test cases
	testCases := make([]*models.TestCase, 0)
	for _, test_case := range test_cases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(test_case), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "unmarshal test case error",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "input is required",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "output is required",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        help.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCases = append(testCases, testCaseBasic)
	}
	data.TestCases = testCases

	//create problem
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create problem error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})
}

// ProblemUpdate
// @Tags         Admin Method
// @Summary      Problem Update
// @Param        Authorization header     string     true  "Authorization"
// @Param        identity  formData     string     true  "identity"
// @Param        title  formData     string     true  "title"
// @Param        content  formData     string     true  "content"
// @Param        max_runtime  formData     int         true  "max_runtime"
// @Param        max_mem  formData     int         true  "max_mem"
// @Param        category_ids  formData     string     false  "category_ids"
// @Param        test_cases  formData     string     true  "test_cases (JSON array, supports multi-line)"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /admin/problem-update [put]
func ProblemUpdate(c *gin.Context) {
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	max_runtime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	max_mem, _ := strconv.Atoi(c.PostForm("max_mem"))
	category_ids := c.PostFormArray("category_ids")
	test_cases := c.PostFormArray("test_cases")
	if identity == "" || title == "" || content == "" || max_runtime == 0 || max_mem == 0 || len(category_ids) == 0 || len(test_cases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "parameters are required",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		//get problem basic
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: max_runtime,
			MaxMem:     max_mem,
		}
		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {

			return err
		}
		//查询问题详情
		err = tx.Where("identity = ?", identity).First(problemBasic).Error
		if err != nil {
			return err
		}

		//关联分类的更新

		//删除关联分类
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		//新增关联分类
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range category_ids {
			category_id, _ := strconv.Atoi(id)
			pcs = append(pcs, &models.ProblemCategory{
				ProblemID:  problemBasic.ID,
				CategoryID: uint(category_id),
			})
		}
		err = tx.Create(pcs).Error
		if err != nil {
			return err
		}
		//更新测试用例
		//删除已存在的测试用例
		err = tx.Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		//新增测试用例
		tcs := make([]*models.TestCase, 0)
		for _, test_case := range test_cases {
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(test_case), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("input format error")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("output format error")
			}
			tcs = append(tcs, &models.TestCase{
				Identity:        help.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "update problem error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "update problem success",
	})
}
