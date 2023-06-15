package gojs

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dop251/goja"
)

var passthroughEnvVars = []string{"HOME", "USER", "USERPROFILE", "TMPDIR", "TMP", "TEMP", "PATH"}

type CallBack = func([]byte)

type Exec struct {
	workDir string            // 工作区
	envs    map[string]string // 环境变量
}

func NewExec() *Exec {
	e := new(Exec)
	e.envs = make(map[string]string)

	// 写入环境变量
	for _, key := range passthroughEnvVars {
		if value := os.Getenv(key); value != "" {
			e.envs[key] = value
		}
	}

	return e
}

// 设置工作区间
func (e *Exec) SetWorkDir(path string) {
	e.workDir = path
}

// 获取工作区间
func (e *Exec) GetWorkDir() string {
	return e.workDir
}

// 设置环境变量
func (e *Exec) AddEnv(k, v string) {
	e.envs[k] = v
}

// 获取环境变量
func (e *Exec) GetEnv(k string) string {
	s, ok := e.envs[k]
	if ok {
		return s
	}

	return ""
}

func (e *Exec) Run(code string, cb CallBack) error {
	ctx := context.Background()
	args := strings.Split(code, " ")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// 设置工作区
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	// 设置环境变量
	cmd.Env = []string{}
	for k, v := range e.envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// 运行
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cb(in.Bytes())
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func NewExecJs(vm *goja.Runtime, e *Exec) *goja.Object {
	o := vm.NewObject()

	o.Set("setWorkDir", func(path string) {
		e.SetWorkDir(path)
	})

	o.Set("getWorkDir", func() string {
		return e.GetWorkDir()
	})

	o.Set("setEnv", func(k, v string) {
		e.AddEnv(k, v)
	})

	o.Set("getEnv", func(k string) string {
		return e.GetEnv(k)
	})

	o.Set("run", func(code string, cb CallBack) string {
		err := e.Run(code, cb)
		if err != nil {
			fmt.Println("run", err.Error())
			return err.Error()
		}

		return ""
	})

	return o
}
