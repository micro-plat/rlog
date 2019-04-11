package elastic

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/utility"

	elastic "gopkg.in/olivere/elastic.v5"
)

const ConfNode = "elastic"
const ConfLogging = "logging"

//Conf elastic配置
type Conf struct {
	Address      string `json:"address" valid:"requrl,required"`
	UserName     string `json:"user-name"`
	Password     string `json:"password"`
	Index        string `json:"index" valid:"ascii"`
	TypeName     string
	WriteTimeout int    `json:"write-timeout" valid:"required"`
	Cron         string `json:"cron" valid:"ascii,required"`
}

//GetConf 获取elastic配置信息
func GetConf(chConf conf.IConf, index string, typeName string) (c *Conf, err error) {
	var chObjConf Conf
	if err = chConf.Unmarshal(&chObjConf); err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(&chObjConf); !b {
		err = fmt.Errorf("elastic search 配置文件有误:%v", err)
		return nil, err
	}
	chObjConf.Index = index
	chObjConf.TypeName = typeName
	return &chObjConf, nil
}

//GetClient 获取elastic client
func GetClient(s component.IContainer, chObjConf *Conf) (c *elastic.Client, err error) {
	esClient, err := elastic.NewClient(elastic.SetURL(chObjConf.Address),
		elastic.SetBasicAuth(chObjConf.UserName, chObjConf.Password))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if exists, err := esClient.IndexExists(chObjConf.Index).Do(ctx); exists || err != nil {
		return esClient, err
	}
	createIndex, err := esClient.CreateIndex(chObjConf.Index).Do(ctx)
	if err != nil {
		err = fmt.Errorf("创建索引%s失败 %v", chObjConf.Index, err)
		return nil, err
	}
	if !createIndex.Acknowledged {
		err = fmt.Errorf("索引%s创建成功但不可用！", chObjConf.Index)
		return nil, err
	}
	return esClient, nil

}

//BenchAddData 添加数据到elastic
func BenchAddData(client *elastic.Client, typeName string, index string, timeout int, datas [][]byte) (n int, err error) {
	bulkRequest := client.Bulk().Index(index).Type(typeName)
	for _, item := range datas {

		logid := utility.GetGUID()
		data := string(item)
		n += utf8.RuneCount(item)
		indexReq := elastic.NewBulkIndexRequest().Index(index).Type(typeName).Id(logid).Doc(data)
		bulkRequest = bulkRequest.Add(indexReq)
	}

	if bulkRequest.NumberOfActions() != len(datas) {
		err = fmt.Errorf("添加数据与生成的bulk数据条数不匹配，数据 %d 条,bulk %d 条", len(datas), bulkRequest.NumberOfActions())
		return 0, err
	}
	ctx := context.TODO()
	var cannel context.CancelFunc
	if timeout > 0 {
		ctx, cannel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cannel()
	}
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		err = fmt.Errorf("添加bulk数据发生错误：%v", err)
		return 0, err
	}
	if bulkResponse == nil {
		err = fmt.Errorf("bulk返回值bulkResponse为nil")
		return 0, err
	}
	return n, nil
}

//AddData 添加数据到elastic
func AddData(client *elastic.Client, logID string, typeName string, index string, timeout int, data string) (err error) {
	rctx := context.TODO()
	var cannel context.CancelFunc
	if timeout > 0 {
		rctx, cannel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cannel()
	}
	ctx, cannel := context.WithTimeout(rctx, time.Second*time.Duration(timeout))
	defer cannel()
	_, err = client.Index().
		Index(index).
		Type(typeName).
		Id(logID).BodyString(data).
		Refresh("true").
		Do(ctx)
	if err != nil {
		err = fmt.Errorf("添加到elastic发生错误:%v", err)
		return err
	}
	return
}
