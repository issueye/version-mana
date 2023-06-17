package gogit

import (
	"os"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/issueye/version-mana/pkg/utils"
)

var RepoMap = new(sync.Map)

// RepoClone
// 从仓库地址拷贝指定分支到本地指定路径
func RepoClone(path string, option *git.CloneOptions) (*git.Repository, error) {
	value, ok := RepoMap.Load(path)
	if ok {
		return value.(*git.Repository), nil
	}

	// 判断本地是否存在指定的文件夹
	exist, err := utils.PathExists(path)
	if err != nil {
		return nil, err
	}

	// 如果存在文件夹，则直接删除文件夹
	if exist {
		err = os.RemoveAll(path)
		if err != nil {
			return nil, err
		}
	}

	// 克隆代码
	r, err := git.PlainClone(path, false, option)
	if err != nil {
		return nil, err
	}

	// 将对象 存入map中
	RepoMap.Store(path, r)
	return r, nil
}

// 分支信息
type BranchInfo struct {
	Name      string `json:"name"`       // 分支名称
	ShortName string `json:"short_name"` // 短名称
	Hash      string `json:"hash"`       // 分支HASH
	Type      string `json:"type"`       // 类型
}

func GetBranchList(r *git.Repository) ([]*BranchInfo, error) {
	ri, err := r.References()
	if err != nil {
		return nil, err
	}

	list := make([]*BranchInfo, 0)
	err = ri.ForEach(func(r *plumbing.Reference) error {
		branch := new(BranchInfo)
		branch.Name = r.Name().String()
		branch.ShortName = r.Name().Short()
		branch.Hash = r.Hash().String()
		branch.Type = r.Type().String()
		list = append(list, branch)

		return nil
	})

	return list, err
}
