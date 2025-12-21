package service

import (
	"GeekCoding/define"
	"GeekCoding/help"
	"GeekCoding/models"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSubmitList
// @Tags         Public Method
// @Summary      Get Submit List
// @Param        page  query     int     false  "page number, default is 1"
// @Param        size  query     int     false  "size"
// @Param        problem_identity  query     string     false  "problem identity"
// @Param        user_identity  query     string     false  "user identity"
// @Param        status  query     int     false  "submit status"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Submit List Page error: ", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]models.Submit_Basic, 0)

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error

	if err != nil {
		log.Println("Get Submit List error: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get submit list error: " + err.Error(),
		})
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

// SubmitCode
// @Tags         User Private Method
// @Summary      Submit Code
// @Param        authorization  header     string     true  "authorization"
// @Param        problem_identity  query     string     true  "problem identity"
// @Param        code  body     string     true  "code"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /user/submit [post]
func SubmitCode(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code error: " + err.Error(),
		})
		return
	}
	//save code
	path, err := help.SaveCode(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "save code error: " + err.Error(),
		})
		return
	}
	//submit code
	u, exists := c.Get("user")
	if !exists || u == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "user not found",
		})
		return
	}
	userClaim, ok := u.(*help.UserClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "invalid user claims",
		})
		return
	}
	sb := &models.Submit_Basic{
		Identity:        help.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}
	//judging
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "problem not found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "get problem error: " + err.Error(),
			})
		}
		return
	}
	if len(pb.TestCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem has no test cases",
		})
		return
	}

	// 先编译代码，检查编译错误（在 for 循环之前）
	codeDir := filepath.Dir(path)
	codeFile := filepath.Base(path)
	binaryUUID := help.GetUUID()
	binaryPath := filepath.Join(codeDir, binaryUUID)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, codeFile)
	buildCmd.Dir = codeDir
	var buildStderr bytes.Buffer
	buildCmd.Stderr = &buildStderr
	if err := buildCmd.Run(); err != nil {
		// 编译失败，直接返回，不执行测试用例
		sb.Status = 5 // 编译错误
		err = models.DB.Transaction(func(tx *gorm.DB) error {
			err = tx.Create(sb).Error
			if err != nil {
				return errors.New("SubmitBasic Save Error:" + err.Error())
			}
			m := make(map[string]interface{})
			m["submit_num"] = gorm.Expr("submit_num + ?", 1)
			err = tx.Model(new(models.User_Basic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
			if err != nil {
				return errors.New("UserBasic Modify Error:" + err.Error())
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "submit error: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]interface{}{
				"status": sb.Status,
				"msg":    buildStderr.String(),
			},
		})
		return
	}

	//wrong answer
	WA := make(chan int)
	//out of memory
	OOM := make(chan int)

	//pass count
	passCount := 0
	var lock sync.Mutex
	var msg string

	// 获取二进制文件的绝对路径
	var absBinaryPath []byte
	absBinaryPath, err = exec.Command("realpath", binaryPath).Output()
	if err != nil {
		// macOS 使用 readlink
		absBinaryPath, err = exec.Command("readlink", "-f", binaryPath).Output()
		if err != nil {
			// 如果都不可用，使用当前工作目录
			pwd, _ := exec.Command("pwd").Output()
			absBinaryPath = []byte(strings.TrimSpace(string(pwd)) + "/" + binaryPath)
		}
	}
	binaryAbsPath := strings.TrimSpace(string(absBinaryPath))

	for _, testCase := range pb.TestCases {
		// 将 testCase 作为参数传递，避免闭包问题
		go func(tc *models.TestCase) { //执行测试
			// 使用 Docker 运行，设置内存限制（单位：MB）
			// --memory: 内存限制
			// --memory-swap: 内存+交换分区限制（设为相同值禁用swap）
			// --network none: 禁用网络访问
			// --rm: 容器退出后自动删除
			// -v: 挂载二进制文件到容器
			memoryLimit := strconv.Itoa(pb.MaxMem) + "m"
			timeoutSeconds := strconv.Itoa(pb.MaxRuntime/1000 + 1) // 转换为秒，加1秒缓冲

			// 构建 Docker 命令
			// 使用 volume 挂载二进制文件，容器内路径为 /app/runner
			dockerCmd := exec.Command("docker", "run",
				"--rm",
				"--memory="+memoryLimit,
				"--memory-swap="+memoryLimit,
				"--network=none",
				"--stop-timeout="+timeoutSeconds,
				"-v", binaryAbsPath+":/app/runner:ro", // 只读挂载
				"-i", // 保持 stdin 开放
				"golang-code-runner:latest",
				"/app/runner",
			)

			var out, stderr bytes.Buffer
			dockerCmd.Stderr = &stderr
			dockerCmd.Stdout = &out
			stdinPipe, err := dockerCmd.StdinPipe()
			if err != nil {
				log.Println("Docker stdin pipe error:", err)
				WA <- 1
				return
			}

			// 启动 Docker 容器
			if err := dockerCmd.Start(); err != nil {
				log.Println("Docker start error:", err)
				WA <- 1
				return
			}

			// 写入输入数据
			io.WriteString(stdinPipe, tc.Input)
			stdinPipe.Close()

			// 等待执行完成
			err = dockerCmd.Wait()

			// 检查是否超内存（Docker 会在超内存时返回特定退出码）
			if err != nil {
				exitError, ok := err.(*exec.ExitError)
				if ok {
					// 退出码 137 通常表示被 SIGKILL 杀死，可能是超内存
					// Docker 的 OOM killer 会发送 SIGKILL
					if exitError.ExitCode() == 137 {
						msg = "运行超内存"
						OOM <- 1
						return
					}
				}
				// 其他运行时错误视为答案错误
				log.Println("Docker run error:", err, stderr.String())
				WA <- 1
				return
			}

			//答案错误
			if tc.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}

			//答案正确
			lock.Lock()
			passCount++
			lock.Unlock()
		}(testCase)
	}
	select {
	//-1 待判断 1 答案正确 2 答案错误 3 运行超时 4 运行超内存 5 编译错误（已在前面处理）
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCases) {
			sb.Status = 1
			msg = "答案正确"
		} else {
			sb.Status = 3
			msg = "运行超时"
		}

	}

	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		//更新user_basic
		err = tx.Model(new(models.User_Basic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("UserBasic Modify Error:" + err.Error())
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "submit error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}
