package log

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/rlog/modules/logging"
)

type SaveHandler struct {
}

//NewSaveHandler 创建服务
func NewSaveHandler() (u *SaveHandler) {
	return &SaveHandler{}
}

//Handle 保存日志记录
func (u *SaveHandler) Handle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------保存日志----------")
	if err := ctx.Request().Check("plat", "system"); err != nil {
		return err
	}

	body, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	if len(body) <= 2 {
		return errs.NewError(204, "无须处理")
	}
	index := fmt.Sprintf("%s_%s", ctx.Request().GetString("plat"), ctx.Request().GetString("system"))
	logger, err := logging.GetLogging(hydra.C.Container())
	if err != nil {
		return err
	}
	if err = logger.Save(body); err != nil {
		return err
	}
	return "success"
}
