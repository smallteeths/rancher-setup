package server

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rancher-setup/internal/server/middleware"
	appRouter "github.com/rancher-setup/internal/server/router"
)

type Http struct {
	engine    *gin.Engine
	logger    *zap.Logger
	router    *appRouter.Router
	afterFunc func()
	mode      string
	port      string
}

type HttpServer interface {
	GetServerEngine() *gin.Engine
	SetMiddleware(middlewares ...middleware.Interface)
}

type Option interface {
	apply(http *Http)
}

type OptionFunc func(http *Http)

func (f OptionFunc) apply(http *Http) {
	f(http)
}

func New(opts ...Option) *Http {
	httpClass := &Http{}
	for _, opt := range opts {
		opt.apply(httpClass)
	}
	httpClass.defaultOption()
	httpClass.engine = httpClass.setServerEngine()
	httpClass.router = appRouter.New(httpClass.engine)
	return httpClass
}

func (h *Http) SetRouters(routers appRouter.Interface) *Http {
	h.router.AddRouter(routers)
	return h
}

func (h *Http) GetServerEngine() *gin.Engine {
	return h.engine
}

func (h *Http) setServerEngine() (engine *gin.Engine) {
	switch h.mode {
	case gin.DebugMode:
		engine = gin.Default()
	default:
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		engine = gin.New()
		engine.Use(gin.Logger(), middleware.New(h.logger).Handle())
	}
	return
}

func (h *Http) SetMiddleware(middlewares ...middleware.Interface) {
	if len(middlewares) == 0 {
		h.engine.Use(middleware.New(h.logger).Handle())
	} else {
		for _, mid := range middlewares {
			h.engine.Use(mid.Handle())
		}
	}
}

func (h *Http) defaultOption() {
	if h.mode == "" {
		h.mode = gin.DebugMode
	}
}

func (h *Http) Run() {
	srv := http.Server{
		Addr:    h.port,
		Handler: h.engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			if h.logger != nil {
				h.logger.Fatal(fmt.Sprintf("listen: %s\n", err.Error()))
			} else {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
	h.afterExec()
	h.ListenSignal(&srv)
}

func (h *Http) afterExec() {
	if h.afterFunc != nil {
		h.afterFunc()
	}
}

func (h *Http) ListenSignal(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if h.logger != nil {
		h.logger.Info("Shutdown Server!")
	} else {
		log.Println("Shutdown Server!")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		if h.logger != nil {
			h.logger.Fatal("Server Shutdown:" + err.Error())
		} else {
			log.Fatal("Server Shutdown:", err)
		}
	}
	if h.logger != nil {
		h.logger.Info("Server exiting!")
	} else {
		log.Println("Server exiting!")
	}
}

func WithMode(mode string) Option {
	return OptionFunc(func(http *Http) {
		http.mode = mode
	})
}

func WithLogger(logger *zap.Logger) Option {
	return OptionFunc(func(http *Http) {
		http.logger = logger
	})
}

func WithPort(port string) Option {
	return OptionFunc(func(http *Http) {
		http.port = port
	})
}

func WithAfterFunc(afterFunc func()) Option {
	return OptionFunc(func(http *Http) {
		http.afterFunc = afterFunc
	})
}
