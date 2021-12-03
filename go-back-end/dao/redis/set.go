package redis

import (
	"fmt"

	"github.com/spf13/viper"
)

//Set 集合
type Set struct {
	key string
}

//NewOptionSet 配置项的集合
func NewOptionSet(id int64) *Set {
	return &Set{
		key: fmt.Sprintf("module_options::%d", id),
	}
}

//NewParamSet 参数的集合
func NewParamSet(id int64) *Set {
	return &Set{
		key: fmt.Sprintf("module_params::%d", id),
	}
}

//NewModuleSetOfProject 项目关联的模块集合
func NewModuleSetOfProject(projectURLKey string) *Set {
	return &Set{
		key: fmt.Sprintf("%s::project::%s::modules", viper.GetString("redis.key_prefix"), projectURLKey),
	}
}

//NewUserBlackList 用户黑名单
func NewUserBlackList() *Set {
	return &Set{
		key: fmt.Sprintf("%s::black_list::user", viper.GetString("redis.key_prefix")),
	}
}

//NewTokenBlackList token黑名单
func NewTokenBlackList() *Set {
	return &Set{
		key: fmt.Sprintf("%s::black_list::token", viper.GetString("redis.key_prefix")),
	}
}

//IsMember 判断成员是否在集合中
func (s *Set) IsMember(member string) (bool, error) {
	return rdb.SIsMember(s.key, member).Result()
}

//Rem 从集合中移除某个成员
func (s *Set) Rem(member string) (int64, error) {
	return rdb.SRem(s.key, member).Result()
}

//Members 获取集合的所有成员
func (s *Set) Members() ([]string, error) {
	return rdb.SMembers(s.key).Result()
}

//Update 更新配置项集合
func (s *Set) Update(members []string) error {
	//1.清空配置项
	_, err := rdb.Del(s.key).Result()
	if err != nil {
		return err
	}
	//2.添加配置项
	for i := range members {
		_, err := rdb.SAdd(s.key, members[i]).Result()
		if err != nil {
			return err
		}
	}
	return err
}

//Clear 清空所有成员
func (s *Set) Clear() error {
	_, err := rdb.Del(s.key).Result()
	return err
}

//Add 添加成员
func (s *Set) Add(member string) (int64, error) {
	return rdb.SAdd(s.key, member).Result()
}
