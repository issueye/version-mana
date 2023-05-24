package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node
var lock = new(sync.Mutex)

func Init(machineID int64) (err error) {
	snowflake.Epoch = time.Now().UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

// GenID
// 生成 64 位的 雪花 ID
func GenID() int64 {
	lock.Lock()
	defer lock.Unlock()
	return node.Generate().Int64()
}

func init() {
	if err := Init(1); err != nil {
		fmt.Println("Init() failed, err = ", err)
		return
	}
}
