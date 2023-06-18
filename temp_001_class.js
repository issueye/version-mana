var exec = require('vmm/exec')
var time = require('std/time')


class Builder {
    // 构造方法
    constructor() {
        this.cmd = exec.new()

        // 获取工作区
        this.workDir = lichee.workDir
        // 直接删除文件夹
        console.log('删除文件夹：', this.workDir)
        removeExists(this.workDir)
        // 检查文件夹，没有则创建
        console.log('创建文件夹：', this.workDir)
        createNotExists(this.workDir)

        // 设置工作区
        console.log('设置工作区：', this.workDir)
        this.cmd.setWorkDir(this.workDir)
    }

    // 设置工作区
    getWorkDir = () => {
        return lichee.workDir
    }

    // 拉取代码
    clone = () => {
        // 拉取代码
        console.log('拉取代码')
        this.cmd.run('git', ['clone', '-c', 'https.proxy="127.0.0.1:1080"', lichee.repo.RepoUrl], function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })
        // 设置工作区
        console.log('设置工作区：', this.workDir + '\\' + lichee.repo.ProjectName)
        this.cmd.setWorkDir(this.workDir + '\\' + lichee.repo.ProjectName)
    }

    // 获取程序信息
    versionInfo = () => {
        // 获取程序信息
        this.tag = lichee.vers.Tag
        this.version = lichee.vers.Version
        this.branch = lichee.vers.Branch
        this.commit = lichee.vers.CommitHash
        this.nowTime = `${time.nowYear()}-${time.nowMonth()}-${time.nowDay()}T${time.nowHour()}:${time.nowMinute()}:${time.nowSecond()}`
        this.nowDate = time.nowYear() + time.nowMonth() + time.nowDay()

        this.description = '医星排队叫号后台服务'
        this.appName = `${lichee.repo.ServerName}_${version}.${nowDate}_${tag}`
    }

    // 环境变量
    env = () => {
        // 设置环境变量
        console.log('设置环境变量')
        this.cmd.setEnv('GOPROXY', 'https://goproxy.cn')
        this.cmd.setEnv('CGO_ENABLED', '0')
        this.cmd.setEnv('GOARCH', 'amd64')
        this.cmd.setEnv('GOCACHE', lichee.cacheDir)

        this.ldflags = `-s -w -X golang.corp.yxkj.com/orange/task/app/initialize.AppName=${description} `
        ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Branch=\`${branch}\` `
        ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Commit=${commit} `
        ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Date=${nowTime} `
        ldflags += `-X golang.corp.yxkj.com/orange/task/app/initialize.Version=${version}`
    }

    // 更新包、拉取到本地
    pull = () => {
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
    }

    // 构建
    build = () => {
        this.compile("windows")
        this.compile("linux")
    }

    // 编译
    compile = (platform) => {
        switch (platform) {
            case "windows":
                {
                    // 编译 windows
                    cmd.setEnv('GOOS', 'windows')
                    // 编译程序
                    let params = [
                        'build',
                        '-o',
                        `bin/${appName}.exe`,
                        //    '-ldflags',
                        //    ldflags,
                        'main.go'
                    ]

                    return params;
                }
            case "linux":
                {
                    console.log('开始编译程序: linux')
                    // 编译程序
                    let params = [
                        'build',
                        '-o',
                        `bin/${appName}`,
                        //    '-ldflags',
                        //    ldflags,
                        'main.go'
                    ]
                    return params;
                }
        }

        cmd.run('go', params, function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })
    }
}

let build = new Builder()

console.log('build', build.workDir)