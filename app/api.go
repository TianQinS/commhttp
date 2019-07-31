package app

import (
	"net/http"
	"net/http/pprof"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"

	"github.com/TianQinS/commhttp/app/middleware"
	"github.com/TianQinS/commhttp/app/sample"
	"github.com/TianQinS/commhttp/config"
)

// The relevant API functions of cache.
func InitApi(app *iris.Application) {
	var getAPI router.Party
	var setAPI router.Party

	if config.Conf.Debug {
		ppApi := app.Party("/debug")
		ppApi.Get("/pprof", pprofHandler(pprof.Index))
		ppApi.Get("/cmdline", pprofHandler(pprof.Cmdline))
		ppApi.Get("/profile", pprofHandler(pprof.Profile))
		ppApi.Post("/symbol", pprofHandler(pprof.Symbol))
		ppApi.Get("/symbol", pprofHandler(pprof.Symbol))
		ppApi.Get("/trace", pprofHandler(pprof.Trace))
		ppApi.Get("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		ppApi.Get("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		ppApi.Get("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		ppApi.Get("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		ppApi.Get("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		ppApi.Get("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))

		getAPI = app.Party("/api/get")
		setAPI = app.Party("/api/set")
		testApi := app.Party("/test")

		testApi.Get("/crontab/{stime:string min(19)}", sample.TestCrontab)
		testApi.Post("/crontab/{cron:string min(9)}", sample.SetCrontab)
	} else {
		setAPI = app.Party("/api/set", middleware.NewSetAuth())
		getAPI = app.Party("/api/get", middleware.NewGetAuth())
	}

	setAPI.Post("/item/{name:string min(1)}/{key:string min(1)}", sample.SetItem)
	setAPI.Delete("/item/{name:string min(1)}/{key:string min(1)}", sample.DeleteItem)
	setAPI.Post("/cache/{name:string min(1)}/{capacity:int min(1)}/{mode:int}", sample.InitCache)
	setAPI.Delete("/cache/{name:string min(1)}", sample.DeleteCache)
	getAPI.Get("/item/{name:string min(1)}/{key:string min(1)}", sample.GetItem)
}

func pprofHandler(f http.HandlerFunc) context.Handler {
	handler := http.HandlerFunc(f)
	return func(ctx iris.Context) {
		handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	}
}
