package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/pkg/utils"
	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
	*BaseService
}

func NewRepo(db *gorm.DB) *Repo {
	repo := new(Repo)
	repo.Db = db
	repo.BaseService = NewBaseService(db)
	return repo
}

// 创建一个代码仓库信息
func (r *Repo) Create(data *model.CreateRepository) error {
	repo := new(model.Repository)
	idStr := strconv.FormatInt(utils.GenID(), 10)
	repo.ID = utils.Sha1(idStr)
	repo.ProjectName = data.ProjectName
	repo.ServerName = data.ServerName
	repo.RepoUrl = data.RepoUrl
	repo.Code = data.Code
	repo.CreateAt = time.Now().Format("2006-01-02 15:04:05.999")
	return r.Db.Create(repo).Error
}

// 修改代码仓库信息
func (r *Repo) Modify(data *model.ModifyRepository) error {
	repo := new(model.Repository)
	repo.ID = data.ID
	repo.ProjectName = data.ProjectName
	repo.ServerName = data.ServerName
	repo.RepoUrl = data.RepoUrl
	repo.Code = data.Code
	return r.Db.Model(repo).Where("id = ?", repo.ID).Updates(repo).Error
}

// 根据id 查询
func (r *Repo) GetById(id string) (*model.Repository, error) {
	data := new(model.Repository)
	err := r.Db.Model(data).Where("id = ?", id).Find(data).Error
	return data, err
}

// 根据条件查询
func (r *Repo) Query(req *model.QueryRepository) ([]*model.Repository, error) {
	list := make([]*model.Repository, 0)
	err := r.DataFilter(model.Repository{}.TableName(), req, &list, func(db *gorm.DB) (*gorm.DB, error) {
		query := db.Order("id")

		if req.Condition != "" {
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Or("path like ?", fmt.Sprintf("%%%s%%", req.Condition))
		}

		return query, nil
	})
	return list, err
}

// 根据 id 删除
func (r *Repo) DelById(id string) error {
	return r.Db.Where("id = ?", id).Delete(&model.Repository{}).Error
}

/****************************************AppVersion*******************************************/

// 根据id 查询
func (r *Repo) GetVersionById(id string) (*model.AppVersionInfo, error) {
	data := new(model.AppVersionInfo)
	err := r.Db.Model(data).Where("id = ?", id).Find(data).Error
	return data, err
}

// 根据分支名称 查询
func (r *Repo) GetVersionByBranch(branch string) ([]*model.AppVersionInfo, error) {
	data := make([]*model.AppVersionInfo, 0)
	err := r.Db.Model(data).Where("branch = ?", branch).Find(data).Error
	return data, err
}

// 根据分支名称 查询
func (r *Repo) GetVersionByBranchAndNum(branch string, num string) (*model.AppVersionInfo, error) {
	data := new(model.AppVersionInfo)
	err := r.Db.
		Model(data).
		Where("branch = ?", branch).
		Where("version = ?", num).
		Find(data).
		Error

	return data, err
}

// 根据 id 删除
func (r *Repo) DelVersionById(id string) error {
	return r.Db.Where("id = ?", id).Delete(&model.AppVersionInfo{}).Error
}

// 添加版本
func (r *Repo) CreateVersion(data *model.CreateVersion) error {
	info := new(model.AppVersionInfo)
	info.AppName = data.AppName
	info.Tag = data.Tag
	info.Branch = data.Branch
	info.CommitHash = data.CommitHash
	info.Content = data.Content
	info.RepoID = data.RepoID
	info.Version = data.Version
	info.VersionX = data.VersionX
	info.VersionY = data.VersionY
	info.VersionZ = data.VersionZ

	// 生成一个内部的可排序的版号   101001
	version_y := utils.StrPad(strconv.FormatInt(data.VersionY, 10), 2, "0", "LEFT")
	version_z := utils.StrPad(strconv.FormatInt(data.VersionZ, 10), 3, "0", "LEFT")
	iv := strconv.FormatInt(data.VersionX, 10) + version_y + version_z
	i, err := strconv.ParseInt(iv, 10, 64)
	if err != nil {
		return err
	}

	info.InternalVersion = i
	info.CreateAt = time.Now().Format("2006-01-02 15:04:05.999")

	// 生成一个ID
	id := utils.GenID()
	idStr := strconv.FormatInt(id, 10)
	info.ID = utils.Sha1(idStr)

	return r.Db.Create(info).Error
}

// 获取列表
func (r *Repo) VersionList(req *model.QueryVersion) ([]*model.AppVersionInfo, error) {
	list := make([]*model.AppVersionInfo, 0)

	err := r.DataFilter(model.AppVersionInfo{}.TableName(), req, &list, func(db *gorm.DB) (*gorm.DB, error) {
		query := db.Order("id")

		if req.Tag != "" {
			query = query.Where("tag = ?", req.Tag)
		}

		if req.Condition != "" {
			query = query.Where("app_name like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Where("version like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Where("content like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Where("branch like ?", fmt.Sprintf("%%%s%%", req.Condition))
		}

		return query, nil
	})

	return list, err
}

func (r *Repo) GetLastVerNum(repoId string, req *model.QryLastVer) (*model.AppVersionInfo, error) {
	data := new(model.AppVersionInfo)

	q := r.Db.
		Model(data).
		Where("repo_id = ?", repoId).
		Where("branch = ?", req.Branch).
		Where("tag = ?", req.Tag).
		Find(data).Order("internal_version desc")

	// sqlStr := q.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Find(data)
	// })

	// fmt.Println("sqlStr", sqlStr)

	err := q.Find(data).Error

	return data, err
}
