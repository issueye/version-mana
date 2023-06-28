var exec = require('vmm/exec')
var filepath = require('std/path/filepath')

let cmd = exec.new()
let workDir = lichee.workDir

// 设置环境变量
function setEnv() {
	cmd.setEnv('GOPROXY', 'https://goproxy.cn')
    cmd.setEnv('CGO_ENABLED', '0')
    cmd.setEnv('GOCACHE', lichee.cacheDir)
}

// 检查工作区
function checkPath() {
    // 直接删除文件夹
    console.log('删除文件夹：', workDir) 
    removeExists(workDir)
    // 检查文件夹，没有则创建
    console.log('创建文件夹：', workDir)
    createNotExists(workDir)

    // 设置工作区
    console.log('设置工作区：', workDir)
    cmd.setWorkDir(workDir)
}

// 拉取代码
function cloneCode() {
	let data = lichee.vers.Branch.split('/')
    let branchName = data[data.length - 1]
    console.log('branchName ', branchName)

    cmd.run('git', ['clone', '-b', branchName, '-c', 'http.proxy=127.0.0.1:7890', lichee.repo.RepoUrl], function (val) {
        let data = convToString(val, 0)
        console.log(data)
    })
  
    // 设置工作区
  	workDir = workDir + '\\' + lichee.repo.ProjectName
    console.log('设置工作区：', workDir)
    cmd.setWorkDir(workDir)
  	
  	// 切换到指定的commit
  	console.log(`切换到 commit => ${lichee.vers.CommitHash}`)
  	cmd.run('git', ['checkout', lichee.vers.CommitHash], function (val) {
    	let data = convToString(val, 0)
      	console.log(data)
    })
}

// go mod
function goMod() {
	console.log('拉取更新引用模块')
    cmd.runWait('go', ['mod', 'tidy'], 10, function (val) {
        let data = convToString(val, 0)
        console.log(data)
    })

    console.log('将引用模块拉取到本项目 vendor 目录')
    cmd.runWait('go', ['mod', 'vendor'], 10, function (val) {
        let data = convToString(val, 0)
        console.log(data)
    })
}

// 编译
function compile(platform) {
	let goos = 'linux'
    let appName = `bin/${lichee.vers.AppName}`
    let name = lichee.vers.AppName
    let t = 1
    
    if (platform == "windows") {
    	// 编译 windows
        goos = 'windows'
        appName += '.exe'
      	name += '.exe'
      	t = 0
    }

    let params = [ 'build', '-o', appName, 'main.go' ]
    cmd.setEnv('GOOS', goos)
    cmd.run('go', params, function (val) {
        let data = convToString(val, 0)
        console.log(data)
    })
  
  	console.log(`[${platform}] =>  ${appName} 编译完成`)
  	
    moveApp(name, appName)
  	
  	// 保存数据
  	if (lichee.runType == 1) {
    	createRelease(lichee.vers.ID, t)
    }
}

// 移动程序到下载目录
function moveApp(appName, path) {
  	let tmp = filepath.join(rootPath(), 'runtime', 'static', 'app')
	let newPath = filepath.join(tmp, lichee.vers.AppName, appName)
    let oldPath = filepath.join(workDir, path)
    createNotExists(tmp)
    moveFile(oldPath, newPath)
  	console.log(`[${appName}]移动完成`)
}

// 运行
function buildApp() {
  	// 环境变量
  	setEnv()
  	// 检查工作区
  	checkPath()
	// 拉取代码
    cloneCode()
    // 同步依赖
    goMod()
    // 编译 windows
    compile('windows')
    // 编译 linux
    compile('linux')
  
  	console.log('+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++')
  	console.log('++++++++++++++++++++++++++++编译完成++++++++++++++++++++++++++++++')
  	console.log('+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++')
}

buildApp()
