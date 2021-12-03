package memory

import (
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

//ConnectMap 战利品管道
var ConnectMap sync.Map

//HasSendMessageGoroutineMap 用户是否有消息推送协程
var HasSendMessageGoroutineMap sync.Map

//AddConn 添加一个连接
func AddConn(userID int64, conn *websocket.Conn) {
	var connList []*websocket.Conn
	value, ok := ConnectMap.Load(userID)
	if !ok {
		zap.L().Debug("该用户未建立连接")
		connList = make([]*websocket.Conn, 0)
	} else {
		zap.L().Debug("该用户已建立连接")
		connList, _ = value.([]*websocket.Conn)
	}
	connList = append(connList, conn)
	ConnectMap.Store(userID, connList)
}

//GetConnList 获取指定用户的连接列表
func GetConnList(userID int64) []*websocket.Conn {
	value, _ := ConnectMap.Load(userID)
	connList, _ := value.([]*websocket.Conn)
	return connList
}

//RemConn 删除一个连接
func RemConn(userID int64, conn *websocket.Conn) {
	var connList []*websocket.Conn
	value, ok := ConnectMap.Load(userID)
	if !ok {
		return
	}
	connList, ok = value.([]*websocket.Conn)
	if !ok {
		return
	}
	//寻找指定连接的下标，如果没有找到 index为-1
	var index int = -1
	for i := range connList {
		if connList[i] == conn {
			index = i
			break
		}
	}
	//如果找到了,就删除该元素
	if index != -1 {
		zap.L().Debug("找到了对应连接，进行删除...")
		connList = append(connList[:index], connList[index+1:]...)
		ConnectMap.Store(userID, connList)
	}
}
