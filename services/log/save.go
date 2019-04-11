package log

import (
	"fmt"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/logsaver/modules/logging"
)

type SaveHandler struct {
	container component.IContainer
}

//NewSaveHandler 创建服务
func NewSaveHandler(container component.IContainer) (u *SaveHandler) {
	return &SaveHandler{
		container: container,
	}
}

//Handle 保存日志记录
func (u *SaveHandler) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------保存日志----------")
	body, err := ctx.Request.GetBody()
	if err != nil {
		return err
	}
	if len(body) <= 2 {
		ctx.Response.SetStatus(204)
		return nil
	}
	index, exists := ctx.Request.Get("plat")
	if !exists {
		return fmt.Errorf("路由配置有误，未找到参数plat")
	}
	typeName, exists := ctx.Request.Get("system")
	if !exists {
		return fmt.Errorf("路由配置有误，未找到参数system")
	}
	// 未来发布的elasticsearch 6.0.0版本为保持兼容，仍然会支持单index，多type结构，但是作者已不推荐这么设置。在elasticsearch 7.0.0版本必须使用单index,单type，多type结构则会完全移除。
	index = fmt.Sprintf("%s_%s", index, typeName)

	logger, err := logging.Get(u.container, index, index)
	if err != nil {
		return err
	}
	if err = logger.Save(body); err != nil {
		return err
	}
	return "success"
}
