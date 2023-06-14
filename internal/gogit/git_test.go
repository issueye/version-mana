package gogit

import (
	"fmt"
	"testing"

	"github.com/issueye/version-mana/pkg/utils"
)

func Test_CloneRepo(t *testing.T) {
	// 拷贝仓库
	url := "https://golang.corp.yxkj.com/orange/task.git"
	r, err := RepoClone("./repo_task", url, "PDJH-V2.1-DEV-001")
	if err != nil {
		t.Errorf("克隆仓库失败，失败原因：%s", err.Error())
	}

	bi, err := GetBranchList(r)
	if err != nil {
		t.Errorf("获取分支列表失败，失败原因：%s", err.Error())
	}

	for _, bi2 := range bi {
		s := utils.Struct2Json(bi2)
		fmt.Println(s)
	}
}
