package ws

import (
	"net/http"
	"web_app/dao/memory"
	"web_app/dao/redis"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

//NewConnect 接收客户端的连接
func NewConnect(c *gin.Context) (*websocket.Conn, error) {
	// 批准客户端升级协议为ws,和客户端建立连接
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 解决跨域问题
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("upgrader.Upgrade failed", zap.Error(err))
		return nil, err
	}
	return conn, nil
}

//SendMessage 消息推送函数（一个用户只能有一个消息推送协程）
func SendMessage(userID int64) {
	//消息推送协程退出，删除掉标志变量
	defer memory.HasSendMessageGoroutineMap.Delete(userID)
	lootQueue := redis.NewLootQueue(userID)
	for {
		connList := memory.GetConnList(userID)
		if len(connList) == 0 {
			return
		}
		//从队列中取出战利品
		lootStr, err := lootQueue.DeQueue()
		if err != nil {
			zap.L().Error("redis.NewLootQueue.Dequeue failed", zap.Error(err))
			return
		}
		//向该用户建立的所有的连接都推送消息
		for _, conn := range connList {
			err = conn.WriteMessage(websocket.TextMessage, []byte(lootStr))
			if err != nil {
				//说明取出了消息，但是没有推成功，应该将该消息放回去
				zap.L().Error("conn.WriteMessage failed", zap.Error(err))
				if err := lootQueue.Recover(lootStr); err != nil {
					zap.L().Error("lootQueue.Recover failed", zap.Error(err))
					return
				}
				return
			}
		}
	}
}

//CheckAlive 检查客户端是否关闭了连接（客户端是否下线）
func CheckAlive(userID int64, c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			zap.L().Error("客户端关闭了连接")
			//服务端关闭连接
			c.Close()
			//删除该用户维护的连接数组中的对应项
			memory.RemConn(userID, c)
			break
		}
	}
}
