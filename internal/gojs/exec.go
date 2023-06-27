package gojs

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

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

func (e *Exec) Run(name string, args []string, cb CallBack) error {
	fmt.Println(name, args)

	cmd := exec.Command(name, args...)

	// 设置工作区
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	// 设置环境变量
	cmd.Env = []string{}
	for k, v := range e.envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errInfo := fmt.Errorf("run %s: %s", err.Error(), stderr.String())
		cb([]byte(errInfo.Error()))
		return errInfo
	}

	return nil
}

func (e *Exec) RunWait(name string, args []string, seconds int, cb CallBack) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(seconds))
	defer cancel()
	fmt.Println(name, args)
	cmd := exec.CommandContext(ctx, name, args...)

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
		errInfo := fmt.Errorf("stdoutPipe %s", err.Error())
		cb([]byte(errInfo.Error()))
		return errInfo
	}

	err = cmd.Start()
	if err != nil {
		errInfo := fmt.Errorf("start %s", err.Error())
		cb([]byte(errInfo.Error()))
		return errInfo
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cb(in.Bytes())
	}

	err = cmd.Wait()
	if err != nil {
		errInfo := fmt.Errorf("wait %s", err.Error())
		cb([]byte(errInfo.Error()))
		return errInfo
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

	o.Set("run", func(name string, args []string, cb CallBack) string {
		err := e.Run(name, args, cb)
		if err != nil {
			fmt.Println("run", err.Error())
			return err.Error()
		}

		return ""
	})

	o.Set("runWait", func(name string, args []string, t int, cb CallBack) string {
		err := e.RunWait(name, args, t, cb)
		if err != nil {
			fmt.Println("runWait err:", err.Error())
			return err.Error()
		}

		return fmt.Sprintln(name, args)
	})

	return o
}
