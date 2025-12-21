package service

import (
	"GeekCoding/define"
	"GeekCoding/help"
	"GeekCoding/models"
	"bytes"
	"context"
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
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
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
	//time limit exceeded
	TLE := make(chan int)
	//accepted (all test cases passed)
	AC := make(chan int)

	//pass count
	passCount := 0
	totalTests := len(pb.TestCases)
	completedCount := 0 // 完成的测试用例数量（包括通过、失败、超时等）
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
			// 创建带超时的 context，每个测试用例独立超时
			// 每个测试用例有独立的超时时间
			testTimeout := time.Millisecond * time.Duration(pb.MaxRuntime)
			ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
			defer cancel()

			// 标记测试用例开始执行
			defer func() {
				lock.Lock()
				completedCount++
				lock.Unlock()
			}()

			// 使用 Docker 运行，设置完整的安全沙盒配置
			// --memory: 内存限制
			// --memory-swap: 内存+交换分区限制（设为相同值禁用swap）
			// --network none: 禁用网络访问
			// --rm: 容器退出后自动删除
			// --cpus: CPU 限制
			// --pids-limit: 进程数限制
			// --read-only: 只读根文件系统
			// --tmpfs: 临时文件系统（如果需要）
			// --user: 非 root 用户
			// --cap-drop: 删除所有 capabilities
			// --security-opt: 安全选项
			// --ulimit: 资源限制
			// -v: 挂载二进制文件到容器（只读）
			memoryLimit := strconv.Itoa(pb.MaxMem) + "m"
			timeoutSeconds := strconv.Itoa(pb.MaxRuntime/1000 + 1) // 转换为秒，加1秒缓冲

			// 构建 Docker 命令 - 完整的安全沙盒配置
			dockerArgs := []string{
				"run",
				"--rm",                                      // 自动删除容器
				"--network=none",                            // 网络隔离
				"--memory=" + memoryLimit,                   // 内存限制
				"--memory-swap=" + memoryLimit,              // 禁用 swap
				"--cpus=1.0",                                // CPU 限制（1 核）
				"--pids-limit=10",                           // 进程数限制（防止 fork bomb）
				"--read-only",                               // 只读根文件系统
				"--tmpfs", "/tmp:rw,noexec,nosuid,size=10m", // 临时文件系统（如果需要）
				"--user=1000:1000",                      // 非 root 用户（sandbox:sandbox）
				"--cap-drop=ALL",                        // 删除所有 capabilities
				"--security-opt=no-new-privileges:true", // 禁止获取新权限
				"--ulimit=nofile=64:64",                 // 文件描述符限制
				"--ulimit=nproc=10:10",                  // 进程数限制
				"--ulimit=stack=8192:8192",              // 栈大小限制
				"--stop-timeout=" + timeoutSeconds,      // 停止超时
				"-v", binaryAbsPath + ":/app/runner:ro", // 只读挂载二进制文件
				"-i", // 保持 stdin 开放
				"golang-code-runner:latest",
				"/app/runner",
			}

			dockerCmd := exec.CommandContext(ctx, "docker", dockerArgs...)

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

			// 等待执行完成或超时
			err = dockerCmd.Wait()

			// 检查是否超时（必须在 Wait() 之后检查，因为 context 超时会自动取消命令）
			// 注意：exec.CommandContext 会在 context 超时时自动杀死进程
			if ctx.Err() == context.DeadlineExceeded {
				// Context 超时，CommandContext 应该已经杀死了进程
				// 但为了确保，再次尝试杀死
				if dockerCmd.Process != nil {
					dockerCmd.Process.Kill()
				}
				// 这个测试用例超时，发送 TLE 信号
				// 注意：如果有多个测试用例，只有第一个超时的会发送信号
				select {
				case TLE <- 1:
					msg = "运行超时"
				default:
					// 如果 TLE channel 已经有信号，不重复发送
				}
				return
			}

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
					// 退出码 124 通常表示被 timeout 命令杀死（超时）
					if exitError.ExitCode() == 124 {
						msg = "运行超时"
						TLE <- 1
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
			// 检查是否所有测试用例都通过了
			if passCount == totalTests {
				// 所有测试用例都通过，发送 AC 信号
				select {
				case AC <- 1:
					// AC 信号已发送
				default:
					// 如果 AC channel 已关闭或已发送，忽略
				}
			}
			lock.Unlock()
		}(testCase)
	}
	// 使用全局超时作为兜底机制
	// 由于测试用例是并发执行的，全局超时应该是：单个测试用例超时 + 缓冲时间
	// 而不是乘以测试用例数量（因为它们是并发执行的）
	// 缓冲时间设置为单个测试用例超时的 50%，确保有足够时间处理所有结果
	bufferTime := time.Millisecond * time.Duration(pb.MaxRuntime) / 2
	globalTimeout := time.After(time.Millisecond*time.Duration(pb.MaxRuntime) + bufferTime)

	select {
	//-1 待判断 1 答案正确 2 答案错误 3 运行超时 4 运行超内存 5 编译错误（已在前面处理）
	case <-AC:
		// 所有测试用例都通过，明确标记为 AC
		sb.Status = 1
		msg = "答案正确"
	case <-WA:
		// 答案错误（优先级高于超时）
		sb.Status = 2
	case <-OOM:
		// 运行超内存
		sb.Status = 4
	case <-TLE:
		// 运行超时（单个测试用例超时）
		sb.Status = 3
	case <-globalTimeout:
		// 全局超时（兜底机制），检查所有测试用例的状态
		lock.Lock()
		// 检查是否所有测试用例都完成了
		if completedCount == totalTests {
			// 所有测试用例都完成了
			if passCount == totalTests {
				// 所有测试用例都通过了，但可能因为竞态条件没有收到 AC 信号
				sb.Status = 1
				msg = "答案正确"
			} else {
				// 有测试用例失败，但已经通过其他 channel 处理了（WA/OOM/TLE）
				// 这里不应该到达，因为其他 channel 应该已经处理了
				// 如果到达这里，说明有测试用例失败但没有发送信号，标记为错误
				sb.Status = 2
				msg = "答案错误"
			}
		} else {
			// 还有测试用例未完成，说明确实超时了
			// 可能是所有测试用例都超时，或者部分超时
			sb.Status = 3
			msg = "运行超时"
		}
		lock.Unlock()
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
