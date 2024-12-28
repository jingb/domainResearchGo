package handler

import (
	"context"
	"net/http"
)

// Handler 处理器接口
type Handler interface {
	Handle(ctx context.Context, req *http.Request) (interface{}, error)
}

// BaseHandler 基础处理器
type BaseHandler struct {
	// 可以添加一些共用的依赖，比如日志、配置等
}
