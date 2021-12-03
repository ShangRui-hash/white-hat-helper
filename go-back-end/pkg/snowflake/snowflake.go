package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

//初始化全局的node节点
//startTime 起始时间，格式如同：2006-01-02
//machineID 机器ID
func Init(startTime string, machineID int64) error {
	//1.指定起始时间
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 1000000
	//2.生成node
	node, err = sf.NewNode(machineID)
	return err
}

//雪花算法生成id
func GenID() int64 {
	return node.Generate().Int64()
}
