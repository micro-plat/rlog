package logging

import (
	"fmt"

	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/logger"
)

const typeNode = "elastic"
const nameNode = "logging"

//GetLogging 获取日志组件
func GetLogging(c container.IContainer) (*LoggingService, error) {
	v, err := c.GetOrCreate(typeNode, nameNode, func(conf *conf.RawConf) (interface{}, error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", typeNode, nameNode)
		}
		c := &Conf{}
		if err := conf.ToStruct(c); err != nil {
			return nil, err
		}
		client, err := GetClient(c)
		if err != nil {
			return nil, err
		}
		return NewLoggingService(client, c, logger.New("logging"))
	})
	if err != nil {
		return nil, err
	}
	return v.(*LoggingService), nil

}
