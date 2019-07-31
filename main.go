// You should move this file to GOPATH and modify the configuration file.
package main

import (
	"flag"

	"github.com/TianQinS/commhttp/app"
	"github.com/kataras/iris"
)

var (
	port = flag.String("p", "23456", "iris http port")
)

func InitApi(app *iris.Application) {
}

func main() {
	// runtime.GOMAXPROCS(2)
	// debug.SetGCPercent(300)

	flag.Parse()
	http := app.NewApp()
	InitApi(http)
	http.Run(iris.Addr(":" + *port))
}
