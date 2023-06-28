package logic

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/gogit"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
	"github.com/issueye/version-mana/pkg/utils"
)

type RepoLogic struct{}

func NewRepo() *RepoLogic {
	return &RepoLogic{}
}

// CreateVersion
// 创建版本
func (RepoLogic) CreateVersion(data *model.CreateVersion) error {

	// 检查当前分支版本号是否已经存在
	versionInfo, err := service.NewRepo(global.DB).GetVersionByBranchAndNum(data.Branch, data.Version)
	if err != nil {
		return fmt.Errorf("获取分支信息失败，失败原因：%s", err.Error())
	}

	if versionInfo != nil {
		if versionInfo.ID != "" {
			return fmt.Errorf("当前分支已经存在[%s]这个版本，请检查", data.Version)
		}
	}

	// 检查并且判断程序名称、版本名称、版号是否符合规范

	// 创建版本
	return service.NewRepo(global.DB).CreateVersion(data)
}

// GetVersionList
// 获取版本列表
func (RepoLogic) GetVersionList(req *model.QueryVersion) ([]*model.AppVersionInfo, error) {
	list, err := service.NewRepo(global.DB).VersionList(req)
	if err != nil {
		return nil, fmt.Errorf("获取版本列表失败，失败原因：%s", err.Error())
	}

	return list, nil
}

// RemoveVersion
// 移除版本
func (RepoLogic) RemoveVersion(id string) error {

	// 查询代码仓库信息
	info, err := service.NewRepo(global.DB).GetVersionById(id)
	if err != nil {
		return fmt.Errorf("查询版本信息失败，失败原因：%s", err.Error())
	}

	if info.ID == "" {
		return errors.New("未查询到版本信息")
	}

	err = service.NewRepo(global.DB).DelVersionById(id)
	if err != nil {
		return fmt.Errorf("移除版本失败，失败原因：%s", err.Error())
	}

	return service.NewRepo(global.DB).DelVersionById(id)
}

func (RepoLogic) BranchList(id string, refresh bool) ([]*gogit.BranchInfo, error) {

	// 查询仓库的地址
	repo, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		return nil, err
	}

	if repo.ID == "" {
		return nil, errors.New("未查询到仓库信息")
	}

	s := filepath.Join("runtime", "git_repo", repo.ServerName, "temp")

	options := &git.CloneOptions{
		URL:               repo.RepoUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if repo.ProxyUrl != "" {
		options.ProxyOptions.URL = repo.ProxyUrl
		options.ProxyOptions.Username = repo.ProxyUser
		options.ProxyOptions.Password = repo.ProxyPwd
	}

	var (
		r *git.Repository
	)

	value, ok := gogit.RepoMap.Load(repo.ID)
	if ok && !refresh {
		r = value.(*git.Repository)
	} else {
		r, err = gogit.RepoClone(repo.ID, s, options)
		if err != nil {
			return nil, err
		}
	}

	return gogit.GetBranchList(r)
}

func (RepoLogic) CreateRelease(versionId string, t int) error {
	avi, err := service.NewRepo(global.DB).GetVersionById(versionId)
	if err != nil {
		return err
	}

	down := new(model.ReleaseInfo)
	down.ID = strconv.FormatInt(utils.GenID(), 10)
	down.RepoID = avi.RepoID
	down.VersionID = avi.ID
	down.AppName = avi.AppName
	down.Tag = avi.Tag
	down.Version = avi.Version
	down.InternalVersion = avi.InternalVersion
	down.Branch = avi.Branch
	down.CommitHash = avi.CommitHash

	down.DownUrl = fmt.Sprintf("/www/app/%s/%s", avi.AppName, avi.AppName)
	if t == 0 {
		down.DownUrl += ".exe"
	}

	down.DownCount = 0
	down.Platform = t
	down.CreateAt = utils.GetNowStr()

	return service.NewRepo(global.DB).CreateRelease(down)
}

func (RepoLogic) ReleaseInc(id string) error {
	return service.NewRepo(global.DB).DownCountInc(id)
}
