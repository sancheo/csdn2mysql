package ids

import (
	"fmt"

	"csdn2wordpress/utils/snowflake"
)

var ids = make([]uint64, 0)

// generateIds 一次生成十个
func generateIds() {
	fmt.Println("Generate ten ids...")
	machineNode := 1
	node, err := snowflake.NewNode(int64(machineNode))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		ids = append(ids, uint64(node.Generate().Int64()))
	}
}

// GenerateID 从切片获取
func GenerateID() (id uint64) {
	if len(ids) < 1 {
		generateIds()
	}
	id = ids[0]
	// 取出后从切片删除
	if len(ids) == 1 {
		ids = ids[:0]
	} else {
		ids = ids[1:]
	}
	return
}
