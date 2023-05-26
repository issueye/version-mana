package model

import (
	"sync"

	git "github.com/go-git/go-git/v5"
)

type Repository struct {
	ID         string `gorm:"column:id;primaryKey;type:nvarchar(100);" json:"id"`       // 编码
	Name       string `gorm:"column:name;type:nvarchar(100);" json:"name"`              // 名称
	IsClone    bool   `gorm:"column:is_clone;type:bit;" json:"isClone"`                 // 代码已经拉取
	Path       string `gorm:"column:path;type:nvarchar(400);" json:"path"`              // 代码仓库路径
	CreateTime string `gorm:"column:create_time;type:nvarchar(100);" json:"createTime"` // 创建时间
}

func (Repository) TableName() string {
	return "repository"
}

type CreateRepository struct {
	Name string `json:"name" binding:"required" label:"名称"`     // 名称
	Path string `json:"path" binding:"required" label:"代码仓库路径"` // 代码仓库路径
}

type OpenRepository struct {
	Name string `json:"name" binding:"required" label:"名称"`     // 名称
	Path string `json:"path" binding:"required" label:"代码仓库路径"` // 代码仓库路径
}

type ModifyRepository struct {
	ID   string `json:"id" binding:"required" label:"编码"` // 编码
	Name string `json:"name" label:"仓库名称"`                // 名称
	Path string `json:"path" label:"代码仓库路径"`              // 代码仓库路径
}

// 查询条件
type QueryRepository struct {
	Condition string `json:"condition" form:"condition"` // 条件
	Page
}

type Branch struct {
	Name   string `json:"name"`   // 分支名称
	Type   string `json:"type"`   // 类型
	Target string `json:"target"` //
	Hash   string `json:"hash"`   //
}

type RepoMap struct {
	sync.Map
}

func (r *RepoMap) Get(id string) (*git.Repository, bool) {
	value, ok := r.Load(id)
	if ok {
		return value.(*git.Repository), true
	}

	return nil, false
}

func (r *RepoMap) Put(id string, value *git.Repository) {
	r.Store(id, value)
}
