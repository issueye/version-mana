package model

type Repository struct {
	ID         string `gorm:"column:id;primaryKey;type:nvarchar(100);" json:"id"`       // 编码
	Name       string `gorm:"column:name;type:nvarchar(100);" json:"name"`              // 名称
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
