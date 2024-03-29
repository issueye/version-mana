package model

type Repository struct {
	ID          string `gorm:"column:id;primaryKey;type:nvarchar(100);comment:编码;" json:"id"`            // 编码
	ProjectName string `gorm:"column:project_name;type:nvarchar(100);comment:项目名称;" json:"project_name"` // 名称
	ServerName  string `gorm:"column:server_name;type:nvarchar(100);comment:服务名称;" json:"server_name"`   // 名称
	RepoUrl     string `gorm:"column:repo_url;type:nvarchar(400);comment:代码仓库地址;" json:"repo_url"`       // 代码仓库路径
	Code        string `gorm:"column:code;type:nvarchar(8000);comment:编译代码;" json:"code"`                // 代码
	ProxyUrl    string `gorm:"column:proxy_url;type:nvarchar(400);comment:代理地址;" json:"proxy_url"`       // 代理
	ProxyUser   string `gorm:"column:proxy_user;type:nvarchar(100);comment:代理用户名;" json:"proxy_user"`    // 代理用户名
	ProxyPwd    string `gorm:"column:proxy_pwd;type:nvarchar(100);comment:代理密码;" json:"proxy_pwd"`       // 代理密码
	CreateAt    string `gorm:"column:create_at;type:nvarchar(100);comment:创建时间;" json:"create_at"`       // 创建时间
}

func (Repository) TableName() string {
	return "repository"
}

type CreateRepository struct {
	ProjectName string `json:"project_name" binding:"required" label:"项目名称"` // 项目名称
	ServerName  string `json:"server_name" binding:"required" label:"项目名称"`  // 服务名称
	Code        string `json:"code" label:"代码"`                              // 代码
	RepoUrl     string `json:"repo_url" binding:"required" label:"代码仓库路径"`   // 代码仓库路径
	ProxyUrl    string `json:"proxy_url"  label:"代理地址"`                      // 代理
	ProxyUser   string `json:"proxy_user"  label:"代理用户"`                     // 代理用户名
	ProxyPwd    string `json:"proxy_pwd"  label:"代理密码"`                      // 代理密码
}

type ModifyRepository struct {
	ID          string `json:"id" binding:"required" label:"编码"`             // 编码
	ProjectName string `json:"project_name" binding:"required" label:"项目名称"` // 项目名称
	ServerName  string `json:"server_name" binding:"required" label:"项目名称"`  // 服务名称
	Code        string `json:"code" label:"代码"`                              // 代码
	RepoUrl     string `json:"repo_url" binding:"required" label:"代码仓库路径"`   // 代码仓库路径
	ProxyUrl    string `json:"proxy_url"  label:"代理地址"`                      // 代理
	ProxyUser   string `json:"proxy_user"  label:"代理用户"`                     // 代理用户名
	ProxyPwd    string `json:"proxy_pwd"  label:"代理密码"`                      // 代理密码
}

// 查询条件
type QueryRepository struct {
	Condition string `json:"condition" form:"condition"` // 条件
	Page
}

type AppVersionInfo struct {
	ID              string `gorm:"column:id;primaryKey;type:nvarchar(100);comment:编码;" json:"id"`                  // 编码
	RepoID          string `gorm:"column:repo_id;type:nvarchar(100);comment:仓库ID;" json:"repo_id"`                 // 仓库ID
	AppName         string `gorm:"column:app_name;type:nvarchar(100);comment:程序名称;" json:"app_name"`               // 程序名称 YHLineServer_v2.10.3.20230613_beta
	Tag             string `gorm:"column:tag;type:nvarchar(50);comment:tag;" json:"tag"`                           // tag	 beta
	Version         string `gorm:"column:version;type:nvarchar(100);comment:版本;" json:"version"`                   // 版本    v2.10.3.20230613
	InternalVersion int64  `gorm:"column:internal_version;type:bigint;comment:内部版号;" json:"internal_version"`      // 内部版号
	VersionX        int64  `gorm:"column:version_x;type:int;comment:主版本;" json:"version_x"`                        // 版本 主版本
	VersionY        int64  `gorm:"column:version_y;type:int;comment:功能版本;" json:"version_y"`                       // 版本 功能版本
	VersionZ        int64  `gorm:"column:version_z;type:int;comment:迭代BUG版本;" json:"version_z"`                    // 版本 迭代BUG版本
	CreateAt        string `gorm:"column:create_at;type:nvarchar(100);comment:创建时间;" json:"create_at"`             // 创建时间
	Branch          string `gorm:"column:branch;type:nvarchar(200);comment:分支名称;" json:"branch"`                   // 分支
	CommitHash      string `gorm:"column:commit_hash;type:nvarchar(200);comment:提交的最后一次 hash;" json:"commit_hash"` // 提交的最后一次 hash
	Content         string `gorm:"column:content;type:nvarchar(2000);comment:更新内容;" json:"content"`                // 更新内容
}

func (AppVersionInfo) TableName() string {
	return "app_version_info"
}

type CreateVersion struct {
	RepoID     string `json:"repo_id" binding:"required" label:"仓库ID"`       // 仓库ID
	AppName    string `json:"app_name" binding:"required" label:"程序名称"`      // 程序名称 YHLineServer_v2.10.3.20230613_beta
	Tag        string `json:"tag" binding:"required" label:"tag"`            // tag	 beta
	Version    string `json:"version" binding:"required" label:"版本"`         // 版本    v2.10.3.20230613
	VersionX   int64  `json:"version_x" label:"主版本号"`                        // 版本 主版本
	VersionY   int64  `json:"version_y" binding:"required" label:"功能版本号"`    // 版本 功能版本
	VersionZ   int64  `json:"version_z" binding:"required" label:"迭代版本号"`    // 版本 迭代BUG版本
	Branch     string `json:"branch" binding:"required" label:"分支名称"`        // 分支
	CommitHash string `json:"commit_hash" binding:"required" label:"提交hash"` // 提交的最后一次 hash
	Content    string `json:"content" binding:"required" label:"更新内容"`       // 更新内容
}

// 查询条件
type QueryVersion struct {
	Condition string `json:"condition" form:"condition"`          // 条件
	Tag       string `json:"tag" form:"tag"`                      // tag
	Branch    string `json:"branch" form:"branch"`                // branch
	RepoID    string `json:"repo_id" form:"repo_id" label:"仓库ID"` // 仓库id
	Page
}

// 获取最后版本
type QryLastVer struct {
	Branch string `json:"branch" form:"branch" binding:"required" label:"分支名称"` // 分支
	Tag    string `json:"tag" form:"tag" binding:"required" label:"发布类型"`       // 发布类型
}

type RepoCode struct {
	ID   string `json:"id" binding:"required" label:"编码"`   // 编码
	Code string `json:"code" binding:"required" label:"代码"` // 代码
}

type ReleaseInfo struct {
	ID              string `gorm:"column:id;primaryKey;type:nvarchar(100);comment:编码;" json:"id"`                  // 编码
	RepoID          string `gorm:"column:repo_id;type:nvarchar(100);comment:仓库编码;" json:"repo_id"`                 // 仓库编码
	VersionID       string `gorm:"column:version_id;type:nvarchar(100);comment:版本编码;" json:"version_id"`           // 版本编码
	AppName         string `gorm:"column:app_name;type:nvarchar(100);comment:程序名称;" json:"app_name"`               // 程序名称 YHLineServer_v2.10.3.20230613_beta
	Tag             string `gorm:"column:tag;type:nvarchar(50);comment:tag;" json:"tag"`                           // tag	 beta
	Version         string `gorm:"column:version;type:nvarchar(100);comment:版本;" json:"version"`                   // 版本    v2.10.3.20230613
	InternalVersion int64  `gorm:"column:internal_version;type:bigint;comment:内部版号;" json:"internal_version"`      // 内部版号
	Branch          string `gorm:"column:branch;type:nvarchar(200);comment:分支名称;" json:"branch"`                   // 分支
	CommitHash      string `gorm:"column:commit_hash;type:nvarchar(200);comment:提交的最后一次 hash;" json:"commit_hash"` // 提交的最后一次 hash
	DownUrl         string `gorm:"column:down_url;type:nvarchar(2000);comment:下载地址;" json:"down_url"`              // 下载地址
	DownCount       int64  `gorm:"column:down_count;type:int;comment:下载次数;" json:"down_count"`                     // 下载次数
	Platform        int    `gorm:"column:platform;type:int;comment:平台;" json:"platform"`                           // 平台 0 windows 1 linux
	CreateAt        string `gorm:"column:create_at;type:nvarchar(100);comment:编译时间;" json:"create_at"`             // 编译时间
}

func (ReleaseInfo) TableName() string {
	return "release_info"
}

type QueryRelease struct {
	Condition string `json:"condition" form:"condition"`          // 条件
	Tag       string `json:"tag" form:"tag"`                      // tag
	Branch    string `json:"branch" form:"branch"`                // branch
	RepoID    string `json:"repo_id" form:"repo_id" label:"仓库ID"` // 仓库id
	Page
}
