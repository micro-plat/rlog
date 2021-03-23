package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"

	elastic "github.com/olivere/elastic/v7"
)

//Conf elastic配置
type Conf struct {
	Address      string `json:"address" valid:"requrl,required"`
	UserName     string `json:"user-name"`
	Password     string `json:"password"`
	WriteTimeout int    `json:"write-timeout" valid:"required"`
	Cron         int    `json:"cron" valid:"required"`
}

//GetAddressArry 获取es集群地址列表
func (c *Conf) GetAddressArry() []string {
	if types.IsEmpty(c.Address) {
		return []string{}
	}

	return strings.Split(c.Address, ",")
}

//Client es client
type Client struct {
	*elastic.Client
	index    string
	typeName string
}

//NewClient 获取elastic client
func NewClient(c *Conf, index string, typeName string) (client *Client, err error) {
	addrs := c.GetAddressArry()
	clt, err := elastic.NewClient(elastic.SetURL(addrs...),
		elastic.SetBasicAuth(c.UserName, c.Password))
	if err != nil {
		return nil, err
	}
	client = &Client{Client: clt, index: index, typeName: typeName}
	err = client.CheckIndexType()
	return client, err
}

//CheckIndexType 检查索引等是否存在
func (client *Client) CheckIndexType() error {
	ctx := context.Background()
	if exists, err := client.IndexExists(client.index).Do(ctx); exists || err != nil {
		return err
	}
	createIndex, err := client.CreateIndex(client.index).Do(ctx)
	if err != nil {
		err = fmt.Errorf("创建索引%s失败 %v", client.index, err)
		return err
	}
	if !createIndex.Acknowledged {
		err = fmt.Errorf("索引%s创建成功但不可用！", client.index)
		return err
	}
	return nil
}

//BenchAddData 添加数据到elastic
func (client *Client) BenchAddData(datas [][]byte, timeout int) (n int, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("批量保存数据发生异常：%v", err1)
		}
	}()

	bulkRequest := client.Bulk().Index(client.index)
	for _, item := range datas {
		logid := getLogid(item)
		data := string(item)
		n += utf8.RuneCount(item)
		indexReq := elastic.NewBulkIndexRequest().Index(client.index).Id(logid).Doc(data)
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
func (client *Client) AddData(logID string, timeout int, data string) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("保存数据发生异常：%v", err1)
		}
	}()
	rctx := context.TODO()
	var cannel context.CancelFunc
	if timeout > 0 {
		rctx, cannel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cannel()
	}
	ctx, cannel := context.WithTimeout(rctx, time.Second*time.Duration(timeout))
	defer cannel()
	_, err = client.Index().
		Index(client.index).
		Type(client.typeName).
		Id(logID).BodyString(data).
		Refresh("true").
		Do(ctx)
	if err != nil {
		err = fmt.Errorf("添加到elastic发生错误:%v", err)
		return err
	}
	return
}

func getLogid(data []byte) string {
	gid := utility.GetGUID()
	m := types.XMap{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Sprintf("%d%s", time.Now().UnixNano(), gid[len(gid)-13:])
	}

	timeStr := m.GetString("time")
	if timeStr == "" {
		return fmt.Sprintf("%d%s", time.Now().UnixNano(), gid[len(gid)-13:])
	}

	timeArry := strings.Split(timeStr, ".")
	t, err := time.ParseInLocation("2006/01/02 15:04:05", timeArry[0], time.Local)
	if err != nil {
		return fmt.Sprintf("%d%s", time.Now().UnixNano(), gid[len(gid)-13:])
	}

	logid := fmt.Sprintf("%d", t.Unix())
	if len(timeArry) >= 2 {
		logid = fmt.Sprintf("%s%s", logid, timeArry[1])
	}

	bc := 32 - len(logid)
	if bc == 0 {
		return logid
	}
	if bc < 0 {
		return logid[:32]
	}
	return fmt.Sprintf("%s%s", logid, gid[len(gid)-bc:])
}
