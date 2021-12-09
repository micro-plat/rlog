package main

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/types"
)

var App = hydra.NewApp(
	hydra.WithPlatName("demo"),
	hydra.WithSystemName("test"),
	hydra.WithServerTypes(http.API))

func main() {
	App.Start()
}

func init() {
	hydra.Conf.API("4567")
	App.Micro("/buried/demo", BuriedHandle)

}

func BuriedHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Debug("---------------数据埋点--------------------")

	ctx.Log().Debug("1. 校验必须参数")
	input := &BuriedRequest{}
	if err := ctx.Request().Bind(input); err != nil {
		return err
	}

	ctx.Log().Debug("2. 数据埋点", len(input.Datas))
	rows := make([]map[string]interface{}, len(input.Datas))
	for i := range input.Datas {
		cur := input.Datas[i]
		bytes, _ := json.Marshal(cur)
		rows[i] = map[string]interface{}{
			"server-ip": net.GetLocalIPAddress(),
			"time":      fmt.Sprintf("%s.%09d", cur.TriggerTime, cur.FlowID),
			"level":     types.DecodeString(cur.ResultStatus, 1, "error", "info"),
			"session":   cur.UUID,
			"content":   string(bytes),
		}
	}
	bytes, _ := json.Marshal(rows)
	ctx.Log().Debugf("数据,rows:%v\n", string(bytes))

	resp, err := hydra.C.RPC().GetRegularRPC().Request("/buried/save@192.168.5.108:7011",
		string(bytes),
		rpc.WithHeader("plat", hydra.G.GetPlatName()))
	if err != nil || resp.GetStatus() != 200 || resp.GetResult() != "success" {
		return err
	}

	ctx.Log().Debug("3. 返回结果")
	return "success"
}

//BuriedRequest 埋点请求数据
type Buried struct {
	UUID         string      `json:"uuid" form:"uuid" m2s:"uuid" valid:"required"`                          //跟踪编号
	FlowID       int64       `json:"flow_id" form:"flow_id" m2s:"flow_id" valid:"required"`                 //流水编号
	TriggerTime  string      `json:"trigger_time" form:"trigger_time" m2s:"trigger_time"  valid:"required"` //触发时间
	EventName    string      `json:"event_name" form:"event_name" m2s:"event_name" valid:"required"`        //事件名称
	UserID       string      `json:"user_id" form:"user_id" m2s:"user_id" `                                 //用户编号
	MerchantNo   string      `json:"merchant_no" form:"merchant_no" m2s:"merchant_no"`                      //商户编号
	PagePath     string      `json:"page_path" form:"page_path" m2s:"page_path"`                            //页面地址
	FuncPath     string      `json:"func_path" form:"func_path" m2s:"func_path"`                            //功能地址
	Params       interface{} `json:"params" form:"params" m2s:"params"`                                     //参数
	ResponseInfo interface{} `json:"response_info" form:"response_info" m2s:"response_info"`                //返回信息
	ResultStatus int         `json:"result_status" form:"result_status" m2s:"result_status"`                //结果状态 0.正常结束，1.异常结束
}

type BuriedRequest struct {
	Datas []*Buried `json:"buried_datas" form:"buried_datas" m2s:"buried_datas" valid:"required"`
}
