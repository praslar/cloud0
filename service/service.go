package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/log"
)

type BaseApp struct {
	Config     AppConfig
	Name       string
	Version    string
	Router     *gin.Engine
	HttpServer *http.Server

	listener net.Listener

	initialized            bool
	healthEndpointDisabled bool
}

func NewApp(name, version string) *BaseApp {
	app := &BaseApp{
		Name:                   name,
		Version:                version,
		Router:                 gin.New(),
		HttpServer:             &http.Server{},
		healthEndpointDisabled: false,
	}

	app.HttpServer.Handler = app.Router

	return app
}

func (app *BaseApp) DisableHealthEndpoint() {
	app.healthEndpointDisabled = true
}

func (app *BaseApp) Initialize() error {
	if err := env.Parse(&app.Config); err != nil {
		return err
	}

	app.HttpServer.ReadTimeout = time.Duration(app.Config.ReadTimeout) * time.Second

	// register error handler
	app.Router.Use(
		ginext.ErrorHandler,
		ginext.RequestIDMiddleware,
		ginext.RequestLogMiddleware(app.Name),
	)

	// register routes
	if !app.healthEndpointDisabled {
		app.Router.GET("/", app.HealthHandler)
	}

	app.initialized = true

	return nil
}

func (app *BaseApp) HealthHandler(c *gin.Context) {
	rsp := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{
		Name:    app.Name,
		Version: app.Version,
	}
	c.JSON(http.StatusOK, rsp)
}

func (app *BaseApp) Start(ctx context.Context) {
	l := log.Tag("BaseApp.Start")
	var err error

	if !app.initialized {
		if err = app.Initialize(); err != nil {
			panic(err)
		}
	}

	if app.listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", app.Config.Port)); err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		l.Printf("start listening on %s", app.listener.Addr().String())
		if err := app.HttpServer.Serve(app.listener); err != nil && err != http.ErrServerClosed {
			l.Error(err)
		}
	}()

	wg.Add(1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		defer func() {
			l.Info("shutting down http server ...")
			shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_ = app.HttpServer.Shutdown(shutCtx)
			cancel()
			wg.Done()
		}()

		select {
		case gotSignal, ok := <-signalCh:
			if !ok {
				// channel close
				return
			}
			l.Printf("got signal: %v", gotSignal)
			return
		case <-ctx.Done():
			l.Printf("context is done")
			return
		}
	}()

	go func() {
		l.Printf("start listening debug server on port %d", app.Config.DebugPort)
		_ = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(app.Config.DebugPort), nil)
	}()

	wg.Wait()
}

func (app *BaseApp) Listener() net.Listener {
	return app.listener
}
