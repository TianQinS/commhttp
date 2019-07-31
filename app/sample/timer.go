package sample

import (
	"fmt"

	"github.com/TianQinS/fastapi/timer"
	"github.com/kataras/iris"
)

// Setting up test data.
func SetCrontab(ctx iris.Context) {
	cron := ctx.Params().Get("cron")
	handle := timer.AddCrontab(cron, "test api", func(args ...interface{}) {
		fmt.Println(cron)
	})

	ctx.JSON(map[string]interface{}{
		"ok":   true,
		"data": handle,
	})
}

// For purposes of testing.
func TestCrontab(ctx iris.Context) {
	sTime := ctx.Params().Get("stime")
	err, info := timer.TestCrontab(sTime)
	ctx.JSON(map[string]interface{}{
		"ok":   true,
		"err":  err,
		"data": info,
	})
}
