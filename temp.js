var exec = require('vmm/exec')
var time = require('std/time')

// 获取工作区
let workDir = lichee.workDir

// 创建一个对象
let cmd = exec.new()
// 设置工作区
cmd.setWorkDir(workDir)

// 拉取代码
cmd.run('git', ['clone', lichee.repoInfo.RepoUrl], function (val) {
    let data = convToString(val, 0)
    console.log(data)
})

// 获取程序信息
let tag = 'beta'
let version = 'v.2.1.8'
let branch = ''
let commit = ''
let nowTime = `${time.nowYear()}-${time.nowMonth()}-${time.nowDay()}T${time.nowHour()}:${time.nowMinute()}:${time.nowSecond()}`
let nowDate = time.nowYear() + time.nowMonth() + time.nowDay()

let description = '医星排队叫号后台服务'
let appName = `${lichee.repoInfo.ServerName}_${version}.${nowDate}_${tag}`

let taskWorkDir = workDir + '\\task'
console.log(taskWorkDir)
cmd.setWorkDir(taskWorkDir)

// 设置环境变量
console.log('设置环境变量')
cmd.setEnv('GOPROXY', 'https://goproxy.cn')
cmd.setEnv('CGO_ENABLED', '0')
cmd.setEnv('GOARCH', 'amd64')
cmd.setEnv('GOOS', 'windows')
cmd.setEnv('GOCACHE', workDir + '\\cache')

console.log('获取分支名称')
// 获取分支
cmd.run('git', ['symbolic-ref', '--short', '-q', 'HEAD'], function (val) {
    let data = convToString(val, 0)
    branch = data
})

console.log('获取提交HASH')
// 获取提交HASH
cmd.run('git', ['rev-parse', '--verify', 'HEAD'], function (val) {
    let data = convToString(val, 0)
    commit = data
})

console.log('拉取更新引用模块')
cmd.runWait('go', ['mod', 'tidy'], 10, function () {
    let data = convToString(val, 0)
    console.log(data)
})

console.log('将引用模块拉取到本项目 vendor 目录')
cmd.runWait('go', ['mod', 'vendor'], 10, function () {
    let data = convToString(val, 0)
    console.log(data)
})

let ldflags = `-s -w -X golang.corp.yxkj.com/orange/task/app/initialize.AppName=${description} `
ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Branch=${branch} `
ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Commit=${commit} `
ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Date=${nowTime} `
ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Version=${version}`

// 编译程序
let params = [
    'build',
    '-o',
    `bin/${appName}.exe`,
    '-ldflags',
    ldflags,
    'main.go'
]

console.log('开始编译程序')
cmd.run('go', params, function (val) {
    let data = convToString(val, 0)
    console.log(data)
})

console.log('编译完成')