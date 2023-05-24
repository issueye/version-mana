package service

import (
	"fmt"
	"strconv"

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
	repo.ID = strconv.FormatInt(utils.GenID(), 10)
	repo.Name = data.Name
	repo.Path = data.Path
	return r.Db.Create(repo).Error
}

// 修改代码仓库信息
func (r *Repo) Modify(data *model.ModifyRepository) error {
	repo := new(model.Repository)
	repo.ID = data.ID
	repo.Name = data.Name
	repo.Path = data.Path
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
