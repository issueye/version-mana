package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/model"
	"github.com/issueye/version-mana/internal/service"
)

var Map = new(model.RepoMap)

func RepoClone(id string) (*git.Repository, error) {
	// 查询
	r, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		return nil, err
	}

	p := filepath.Join("git", r.Name)
	repo, err := git.PlainClone(p, false, &git.CloneOptions{
		URL: r.Path,
	})
	if err != nil {
		return nil, err
	}

	// 添加到map中
	Map.Put(r.ID, repo)

	return repo, nil
}

func RepoOpen(id string) (*git.Repository, error) {
	// 查询
	r, err := service.NewRepo(global.DB).GetById(id)
	if err != nil {
		return nil, err
	}

	p := filepath.Join("git", r.Name)
	repo, err := git.PlainOpen(p)
	if err != nil {
		return nil, err
	}

	// 添加到map中
	Map.Put(r.ID, repo)

	return repo, nil
}

func GetBranchList(id string) ([]*model.Branch, error) {
	repo, ok := Map.Get(id)
	if ok {
		refs, err := repo.References()
		if err != nil {
			return nil, err
		}

		head := new(model.Branch)

		list := make([]*model.Branch, 0)
		refs.ForEach(func(r *plumbing.Reference) error {

			fmt.Println("r.String()", r.String())
			// 处理
			branch := new(model.Branch)
			branch.Name = r.Name().Short()
			branch.Hash = r.Hash().String()
			if r.Name().IsRemote() {
				branch.Type = "remote"
			}

			if r.Name().IsBranch() {
				branch.Type = "local"
			}

			if r.Name().String() == "HEAD" {
				head.Name = r.Target().Short()
				head.Hash = r.Hash().String()
				head.Type = r.Name().String()
				return nil
			}

			list = append(list, branch)
			return nil
		})

		data, _ := json.Marshal(head)
		fmt.Println("head", string(data))

		for i, b := range list {
			if b.Name == head.Name {
				head.Hash = b.Hash
				list[i] = head
			}
		}

		return list, nil
	}

	return nil, errors.New("未找到分支列表")
}
