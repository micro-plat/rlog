package services

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
)

//SaveHandle 保存日志记录
func SaveHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------保存日志----------")
	if err := ctx.Request().Check("plat", "system"); err != nil {
		return err
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
	index := fmt.Sprintf("%s_%s", ctx.Request().Headers().GetString("plat"), ctx.Request().Headers().GetString("system"))
	logger, err := GetLogging(hydra.C.Container(), index, index)
	if err != nil {
		return err
	}
	if err = logger.Save(body); err != nil {
		return err
	}
	return "success"
}
