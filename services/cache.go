package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/logger"
)

const typeNode = "elastic"
const nameNode = "logging"

//初始化历史版本处理对象
var historyKeys = newHistories()

//GetLogging 获取日志组件
func GetLogging(c container.IContainer, group, index string, typeName string) (*Logging, error) {
	v, err := c.GetOrCreate(typeNode, nameNode, func(conf *conf.RawConf, keys ...string) (interface{}, error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", typeNode, nameNode)
		}
		c := &Conf{}
		if err := conf.ToStruct(c); err != nil {
			return nil, err
		}

		client, err := NewClient(c, keys[0], keys[1])
		if err != nil {
			return nil, err
		}

		logClient, err := NewLogging(client, c, logger.New("logging"))
		if err != nil {
			return nil, err
		}

		historyKeys.Add(group, index)
		return logClient, nil
	}, index, typeName)
	if err != nil {
		return nil, err
	}
	return v.(*Logging), nil

}

//GetClearClient 获取日志清理组件
func GetClearClient(c container.IContainer) (*ClearClient, error) {
	v, err := c.GetOrCreate(typeNode, nameNode, func(conf *conf.RawConf, keys ...string) (interface{}, error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", typeNode, nameNode)
		}
		c := &Conf{}
		if err := conf.ToStruct(c); err != nil {
			return nil, err
		}

		client, err := NewClearClient(c)
		if err != nil {
			return nil, err
		}
		return client, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*ClearClient), nil

}

//history 历史key记录
type history struct {
	current string
	keys    []string
}

//histories 所有配置的历史key记录
type histories struct {
	records map[string]*history
	lock    sync.Mutex
}

func newHistories() *histories {
	his := &histories{
		records: make(map[string]*history),
	}
	go his.clear()
	return his
}

//Add 添加key信息
func (v *histories) Add(group string, key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	his, ok := v.records[group]
	if !ok {
		v.records[group] = &history{current: key, keys: []string{}}
		return
	}

	his.keys = append(his.keys, his.current)
	his.current = key
}

//Remove 移除key信息
func (v *histories) Remove() {
	v.lock.Lock()
	defer v.lock.Unlock()
	var err error
	for _, history := range v.records {
		for _, k := range history.keys {
			if err1 := hydra.C.Container().Remove(typeNode, nameNode, []string{k, k}...); err1 != nil {
				err = err1
			}
		}
		if err == nil {
			history.keys = []string{}
		}
		err = nil
	}
}

func (v *histories) clear() {
	tk := time.NewTicker(6 * time.Hour)
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-tk.C:
			v.Remove()
		}
	}
}
