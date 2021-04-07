package services

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
)

//SaveHandle 保存日志记录
func SaveHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Debug("--------保存日志----------")
	plat := ctx.Request().Headers().GetString("Plat")
	system := ctx.Request().Headers().GetString("System")
	if plat == "" || system == "" {
		return fmt.Errorf("请求头信息plat和system不能为空")
	}

	//获取数据
	body, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	if len(body) <= 2 {
		return errs.NewError(204, "无须处理")
	}

	//保存日志
	index := fmt.Sprintf("%s_%s%s", plat, system, time.Now().Format("20060102"))
	group := fmt.Sprintf("%s_%s", plat, system)
	logger, err := GetLogging(hydra.C.Container(), group, index, index)
	if err != nil {
		return err
	}
	if err = logger.Save(body); err != nil {
		return err
	}
	return "success"
}

//ClearHandle 清理过期日志
func ClearHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Debug("--------清理过期日志----------")
	logger, err := GetClearClient(hydra.C.Container())
	if err != nil {
		return err
	}
	if err = logger.Clear(); err != nil {
		return err
	}

	return "success"
}
