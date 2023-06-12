package model

type Repository struct {
	ID          string `gorm:"column:id;primaryKey;type:nvarchar(100);" json:"id"`          // 编码
	ProjectName string `gorm:"column:project_name;type:nvarchar(100);" json:"project_name"` // 名称
	ServerName  string `gorm:"column:server_name;type:nvarchar(100);" json:"server_name"`   // 名称
	Path        string `gorm:"column:path;type:nvarchar(400);" json:"path"`                 // 代码仓库路径
	CreateAt    string `gorm:"column:create_at;type:nvarchar(100);" json:"create_at"`       // 创建时间
}

func (Repository) TableName() string {
	return "repository"
}

type CreateRepository struct {
	ProjectName string `json:"project_name" binding:"required" label:"项目名称"` // 项目名称
	ServerName  string `json:"server_name" binding:"required" label:"项目名称"`  // 服务名称
	Path        string `json:"path" binding:"required" label:"代码仓库路径"`       // 代码仓库路径
}

type ModifyRepository struct {
	ID          string `json:"id" binding:"required" label:"编码"`             // 编码
	ProjectName string `json:"project_name" binding:"required" label:"项目名称"` // 项目名称
	ServerName  string `json:"server_name" binding:"required" label:"项目名称"`  // 服务名称
	Path        string `json:"path" binding:"required" label:"代码仓库路径"`       // 代码仓库路径
}

// 查询条件
type QueryRepository struct {
	Condition string `json:"condition" form:"condition"` // 条件
	Page
}
