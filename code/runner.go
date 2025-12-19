package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	//go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	// 写入输入数据
	io.WriteString(stdinPipe, "25 11\n")
	// 关闭 stdin，表示输入完成
	stdinPipe.Close()

	//根据测试的书团里进行运行，拿到输出结果和保准的输出结果进行比较，如果一致则返回正确，否则返回错误
	if err := cmd.Wait(); err != nil {
		log.Fatalln(err, stderr.String())
	}
	fmt.Println(out.String())

	println(out.String() == "36\n")

}
