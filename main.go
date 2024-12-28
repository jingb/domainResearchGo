package main

import (
	"domain-analyzer/config"
	"domain-analyzer/internal/handler"
	"domain-analyzer/internal/pkg/logger"
	"domain-analyzer/internal/service/domain"
	"domain-analyzer/internal/service/ocr"
	"log"
	"path/filepath"

	"domain-analyzer/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(filepath.Join("config", "config.json"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.InitLogger()

	// 初始化OCR服务
	ocrService, err := ocr.NewTencentOCR(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize OCR service: %v", err)
	}

	// 初始化WebArchive服务
	domain.InitWebArchive(cfg)

	// 初始化handler
	h := handler.NewUploadHandler(ocrService)

	r := gin.Default()

	// 添加统一的响应处理中间件
	r.Use(func(c *gin.Context) {
		c.Next()

		// 检查是否已经有响应写入
		if c.Writer.Written() {
			return
		}

		// 定义统一的响应结构
		type Response struct {
			Message string      `json:"message,omitempty"`
			Data    interface{} `json:"data,omitempty"`
		}

		// 获取错误信息
		if err, exists := c.Get("handler_error"); exists {
			logger.Errorf("handler error: %+v", err)
			errObj := err.(error)
			resp := Response{
				Message: errObj.Error(),
			}

			if errors.IsClientError(errObj) {
				c.JSON(400, resp)
			} else if errors.IsServerError(errObj) {
				c.JSON(500, resp)
			} else {
				// 未知错误类型按服务端错误处理
				resp.Message = "Internal Server Error"
				c.JSON(500, resp)
			}
			return
		}

		// 获取handler的处理结果
		if data, exists := c.Get("handler_result"); exists {
			c.JSON(200, Response{
				Data: data,
			})
			return
		}

	})

	// 提供静态文件服务
	r.Static("/static", "./web/static")
	// 加载HTML模板
	r.LoadHTMLGlob("web/template/*")

	// 路由设置
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 包装handler接口调用的辅助函数
	wrapHandler := func(h handler.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			result, err := h.Handle(c.Request.Context(), c.Request)
			if err != nil {
				c.Set("handler_error", err)
				return
			}
			c.Set("handler_result", result)
		}
	}

	r.POST("/upload", wrapHandler(h))

	log.Fatal(r.Run(":" + cfg.Server.Port))
}
