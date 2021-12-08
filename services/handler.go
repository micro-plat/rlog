package services

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
)

//SaveHandle 保存日志记录
func SaveHandle(ctx hydra.IContext) (r interface{}) {
	//return "success"
	///*
	ctx.Log().Debug("--------保存日志----------")
	plat := ctx.Request().Headers().GetString("Plat")
	system := ctx.Request().Headers().GetString("System")
	if plat == "" || system == "" {
		err := fmt.Errorf("请求头信息plat和system不能为空")
		ctx.Log().Error(err)
		return
	}

	//获取数据
	body, err := ctx.Request().GetBody()
	if err != nil {
		err = fmt.Errorf("获取请求内容出错：%w,Body:%s", err, string(body))
		ctx.Log().Error(err)
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
		err = fmt.Errorf("获取Logger：%w", err)
		ctx.Log().Error(err)
		return err
	}
	if err = logger.Save(body); err != nil {
		err = fmt.Errorf("保存日志内容：%w", err)
		ctx.Log().Error(err)
		return err
	}
	return "success"
	//*/
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

//BuriedHandle 保存埋点
func BuriedHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Debug("--------保存埋点----------")
	plat := ctx.Request().Headers().GetString("plat")
	if plat == "" {
		plat = ctx.Request().Headers().GetString("Plat")
	}
	system := "buried"
	if plat == "" || system == "" {
		err := fmt.Errorf("请求头信息plat和system不能为空")
		ctx.Log().Error(err, ctx.Request().Headers())
		return
	}

	//获取数据
	body, err := ctx.Request().GetBody()
	if err != nil {
		err = fmt.Errorf("获取请求内容出错：%w,Body:%s", err, string(body))
		ctx.Log().Error(err)
		return err
	}
	ctx.Log().Debug("body:", string(body))
	if len(body) <= 2 {
		return errs.NewError(204, "无须处理")
	}

	//保存日志
	index := fmt.Sprintf("%s_%s%s", plat, system, time.Now().Format("20060102"))
	group := fmt.Sprintf("%s_%s", plat, system)
	logger, err := GetLogging(hydra.C.Container(), group, index, index)
	if err != nil {
		err = fmt.Errorf("获取Logger：%w", err)
		ctx.Log().Error(err)
		return err
	}
	if err = logger.Save(body); err != nil {
		err = fmt.Errorf("保存日志内容：%w", err)
		ctx.Log().Error(err)
		return err
	}
	return "success"
}
