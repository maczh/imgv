package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sadlil/gologger"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image/gif"
	"image/jpeg"
	"image/png"
	"imgv/controller"
	"imgv/service"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var logger = gologger.GetLogger()

//初始化命令行参数
func parseArgs() (int, string) {
	port := 8080
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cacheDir := ""
	flag.IntVar(&port, "p", 8080, "侦听端口号")
	flag.StringVar(&cacheDir, "d", path+"/cache", "图片缓存目录")
	flag.Parse()
	if strings.HasSuffix(cacheDir, "/") {
		cacheDir = strings.TrimSuffix(cacheDir, "/")
	}
	return port, cacheDir
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors())
	engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 Not Found")
	})
	engine.GET("/image/process", func(c *gin.Context) {
		contentType, img, err := controller.ImageProcess(c)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			w := c.Writer
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", contentType)
			switch contentType {
			case "image/png":
				err = png.Encode(w, img)
			case "image/jpeg":
				err = jpeg.Encode(w, img, &jpeg.Options{Quality: 75})
			case "image/gif":
				err = gif.Encode(w, img, nil)
			case "image/bmp":
				err = bmp.Encode(w, img)
			case "image/tiff":
				err = tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
			default:
				w.Header().Set("Content-Type", "image/png")
				err = png.Encode(w, img)
			}
			w.Flush()
		}
		if err != nil {
			logger.Error("image encode error: " + err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
	})

	return engine
}

func main() {
	port, cacheDir := parseArgs()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0666)
	}
	service.CacheDir = cacheDir
	engine := setupRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server listen: " + err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	logger.Error("Get Signal:" + sig.String())
	logger.Error("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown:" + err.Error())
	}
	logger.Error("Server exiting")
}
