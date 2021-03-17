package services

import (
	"context"
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/lib4go/types"
	elastic "github.com/olivere/elastic/v7"
)

//ClearConf 日志清理配置
type ClearConf struct {
	ExpireDay int64 `json:"expireDay"`
}

//GetExpireDay 获取过期时间
func (n *ClearConf) GetExpireDay() int64 {
	if n.ExpireDay <= 0 {
		//不设置 默认保存15天的日志
		return 15
	}

	return n.ExpireDay
}

//ClearClient .
type ClearClient struct {
	*elastic.Client
}

//NewClearClient 获取elastic client
func NewClearClient(c *Conf) (client *ClearClient, err error) {
	clt, err := elastic.NewClient(elastic.SetURL(c.Address),
		elastic.SetBasicAuth(c.UserName, c.Password))
	if err != nil {
		return nil, err
	}
	client = &ClearClient{Client: clt}
	return client, err
}

//Clear 执行清理
func (n *ClearClient) Clear() error {
	//获取所有的index列表
	ctx := context.Background()
	res, err := n.Client.CatIndices().Do(ctx)
	if err != nil {
		return fmt.Errorf("清理日志,获取index列表异常,err:%v", err)
	}

	hydra.CurrentContext().Log().Debugf("查询的日志列表:%d", len(res))
	if res == nil || len(res) <= 0 {
		return nil
	}

	varConf, err := app.Cache.GetVarConf()
	if err != nil {
		return fmt.Errorf("清理日志时无法获取var.conf:%w", err)
	}

	conf := &ClearConf{}
	_, err = varConf.GetObject("conf", "clearConf", conf)
	if err != nil {
		return fmt.Errorf("清理日志时无法获取var.conf.clearConf:%w", err)
	}

	expireHour := time.Duration(conf.GetExpireDay() * 24 * -1)
	var clearList []string
	expireDate := types.GetInt64(time.Now().Add(expireHour * time.Hour).Format("20060102"))
	//对index列表进行格式化
	for _, item := range res {
		if len(item.Index) < 10 {
			hydra.CurrentContext().Log().Debug("该index不是规范,不进行清理,index:", item.Index)
			continue
		}

		date := types.GetInt64(item.Index[len(item.Index)-8:], 0)
		if date < 19000101 {
			hydra.CurrentContext().Log().Debug("该index不是有效时间内的index,不进行清理,index:", item.Index)
			continue
		}

		hydra.CurrentContext().Log().Debugf("date:%d,expireDate:%d", date, expireDate)
		if date <= expireDate {
			clearList = append(clearList, item.Index)
		}
	}

	//进行删除
	for _, index := range clearList {
		res, err := n.Client.DeleteIndex(index).Do(ctx)
		if err != nil || !res.Acknowledged {
			hydra.CurrentContext().Log().Errorf("删除日志异常,index:%s,res:%v,err:%v", index, res, err)
			continue
		}
	}

	return nil
}
