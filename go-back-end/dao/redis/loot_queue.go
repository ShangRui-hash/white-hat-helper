package redis

import (
	"fmt"
)

//LootQueue 战利品队列 用于消息推送
type LootQueue struct {
	key string
}

//NewLootQueue 构造函数
func NewLootQueue(userID int64) *LootQueue {
	return &LootQueue{
		key: fmt.Sprintf("LootQueue::%d", userID),
	}
}

//EnQueue 入队
func (l LootQueue) EnQueue(lootStr string) error {
	_, err := rdb.LPush(l.key, lootStr).Result()
	return err
}

//DeQueue 出队
func (l LootQueue) DeQueue() (lootStr string, err error) {
	lootPair, err := rdb.BRPop(0, l.key).Result()
	return lootPair[1], err
}

//Recover 恢复，适用于从队列中取出了数据，但是又没有推送成功，这时需要将数据再放回去
func (l LootQueue) Recover(lootStr string) error {
	_, err := rdb.RPush(l.key, lootStr).Result()
	return err
}
