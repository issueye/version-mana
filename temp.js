var exec = require('vmm/exec')
var filepath = require('std/path/filepath')
var time = require('std/time')

let cmd = exec.new()
let workDir = lichee.workDir

// 设置环境变量
function setEnv() {
  console.log('设置编译环境变量')
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

  let errInfo = cmd.run('git', ['clone', '-b', branchName, '-c', 'http.proxy=127.0.0.1:1080', lichee.repo.RepoUrl], function (val) {
    let data = convToString(val, 0)
    console.log(data)
  })

  if (errInfo != '') {
    throw errInfo;
  }

  // 设置工作区
  workDir = workDir + '\\' + lichee.repo.ProjectName
  console.log('设置工作区：', workDir)
  cmd.setWorkDir(workDir)

  // 切换到指定的commit
  console.log(`切换到指定提交： commit => ${lichee.vers.CommitHash}`)
  errInfo = cmd.run('git', ['checkout', lichee.vers.CommitHash], function (val) {
    let data = convToString(val, 0)
    console.log(data)
  })

  if (errInfo != '') {
    throw errInfo;
  }
}

// go mod
function goMod() {
  console.log('拉取更新引用模块 --> 开始')
  let errInfo = cmd.run('go', ['mod', 'tidy'], function (val) {
    let data = convToString(val, 0)
    console.log(data)
  })

  if (errInfo != '') {
    throw errInfo;
  }

  console.log('拉取更新引用模块 --> 成功')
  console.log('创建[vendor] --> 开始')
  errInfo = cmd.run('go', ['mod', 'vendor'], function (val) {
    let data = convToString(val, 0)
    console.log(data)
  })

  if (errInfo != '') {
    throw errInfo;
  }

  console.log('创建[vendor] --> 成功')
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

  let nowDate = time.nowDate()

  let _appName = `-X github.com/issueye/version-mana/internal/initialize.AppName=${name}`
  let _branch = `-X github.com/issueye/version-mana/internal/initialize.Branch=${lichee.vers.Branch}`
  let _commit = `-X github.com/issueye/version-mana/internal/initialize.Commit=${lichee.vers.CommitHash}`
  let _date = `-X github.com/issueye/version-mana/internal/initialize.Date=${nowDate}`
  let _version = `-X github.com/issueye/version-mana/internal/initialize.Version=${lichee.vers.Version}`

  let params = [
    'build',
    `-o=${appName}`,
    '-tags=ui',
    `-ldflags=-w -s ${_appName} ${_branch} ${_commit} ${_date} ${_version}`,
    '.'
  ]
  cmd.setEnv('GOOS', goos)
  console.log(`编译[${platform}]平台程序 --> 开始`)
  let errInfo = cmd.run('go', params, function (val) {
    let data = convToString(val, 0)
    console.log(data)
  })

  if (errInfo != '') {
    throw errInfo;
  }

  console.log(`编译[${platform}]平台程序 --> 编译成功 ^^^ 程序名称：${appName}`)
  moveApp(name, appName, platform)

  // 保存数据
  if (lichee.runType == 1) {
    createRelease(lichee.vers.ID, t)
  }
}

// 移动程序到下载目录
function moveApp(appName, path, platform) {
  let tmp = filepath.join(rootPath(), 'runtime', 'static', 'app')
  let newPath = filepath.join(tmp, lichee.vers.AppName, appName)
  let oldPath = filepath.join(workDir, path)
  let descName = filepath.join(tmp, lichee.vers.AppName, `${platform}_README.txt`)

  createNotExists(tmp)
  console.log(`将编译程序[${appName}]移动到下载目录 --> 开始`)
  moveFile(oldPath, newPath)
  console.log(`将编译程序[${appName}]移动到下载目录 --> 成功`)

  // 创建一个文件对象
  openFile(descName, function (w) {
    w.WriteString(`程序名称：${appName}`)
    w.WriteString(`程序版本：${lichee.vers.Version}`)
    w.WriteString(`最后HASH：${lichee.vers.CommitHash}`)
    w.WriteString(`编译日期：${time.nowDate()}`)
    w.WriteString(`编译分支：${lichee.vers.Branch}`)

    let md5 = fileMD5(newPath)
    w.WriteString(`MD5：${md5}`)

    let sha256 = fileSha256(newPath)
    w.WriteString(`SHA256:${sha256}`)

    w.WriteString(`更新内容：`)
    w.WriteString(lichee.vers.Content)
  })
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
