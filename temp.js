var exec = require('vmm/exec')

class Builder {
  	// 构造函数
  	constructor() {
      	// 命令行执行对象
    	this.cmd = exec.new()
      	
      	this.cmd.setEnv('GOPROXY', 'https://goproxy.cn')
        this.cmd.setEnv('CGO_ENABLED', '0')
      	this.cmd.setEnv('GOCACHE', lichee.cacheDir)
      	
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
  	
  	// 拉取代码
    clone = () => {
        // 拉取代码
        console.log('拉取代码')
      	let data = lichee.vers.Branch.split('/')
        let branchName = data[data.length - 1]
        console.log('branchName ', branchName)
      	
        this.cmd.run('git', ['clone', '-b', branchName, '-c', 'http.proxy=192.168.227.120:7890', lichee.repo.RepoUrl], function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })
        // 设置工作区
        console.log('设置工作区：', this.workDir + '\\' + lichee.repo.ProjectName)
        this.cmd.setWorkDir(this.workDir + '\\' + lichee.repo.ProjectName)
    }
  
  	// 更新包、拉取到本地
    pull = () => {
        console.log('拉取更新引用模块')
        this.cmd.runWait('go', ['mod', 'tidy'], 10, function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })

        console.log('将引用模块拉取到本项目 vendor 目录')
        this.cmd.runWait('go', ['mod', 'vendor'], 10, function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })
    }
  
  	// 编译
    compile = (platform) => {
      	let goos = ''
        let appName = `bin/${lichee.vers.AppName}`
        switch (platform) {
            case "windows":
                {
                    // 编译 windows
                  	goos = 'windows'
                  	appName += '.exe'
                    break;
                }
            case "linux":
                {
                  	goos = 'linux'
                    break;
                }
        }
		
      	let params = [ 'build', '-o', appName, 'main.go' ]
      	this.cmd.setEnv('GOOS', goos)
        this.cmd.run('go', params, function (val) {
            let data = convToString(val, 0)
            console.log(data)
        })
    }

  
  	// 运行构建任务
  	run = () => {
      	// 拉取代码
      	this.clone()
      	// 同步依赖
      	this.pull()
      	// 编译 windows
    	this.compile('windows')
      	// 编译 linux
      	this.compile('linux')
    }
}

let build = new Builder()
build.run()