// Provides basic authentication for specific routes or for the whole app via middleware.
package middleware

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/TianQinS/commhttp/config"
	"github.com/TianQinS/fastapi/timer"
	"github.com/kataras/iris/context"
)

var (
	conf = config.Conf.App
	// staging token data by name for performance.
	tempForbidRemote map[string]int    = make(map[string]int, conf.DefaultMapSize)
	tempClientToken  map[string]string = make(map[string]string, conf.DefaultMapSize)
	tempServerToken  map[string]string = make(map[string]string, conf.DefaultMapSize)
)

func init() {
	// the suspicious address will be banned for thirty minutes.
	timer.AddCrontab("*/30 * * * *", "auth ban list", func() {
		tempForbidRemote = map[string]int{}
	})
}

func getClientToken(key string) string {
	if token, ok := tempClientToken[key]; ok {
		return token
	} else {
		s := sha1.New()
		m := md5.New()
		io.WriteString(s, key)
		io.WriteString(m, fmt.Sprintf(conf.FormatSignSaltClient, s.Sum(nil)))
		token := fmt.Sprintf("%x", m.Sum(nil))
		tempClientToken[key] = token
		return token
	}
	return ""
}

func getServerToken(key string) string {
	if token, ok := tempServerToken[key]; ok {
		return token
	} else {
		s := sha1.New()
		m := md5.New()
		io.WriteString(s, key)
		io.WriteString(m, fmt.Sprintf(conf.FormatSignSaltServ, s.Sum(nil)))
		token := fmt.Sprintf("%x", m.Sum(nil))
		tempServerToken[key] = token
		return token
	}
	return ""
}

// NewGetAuth returns a new auth middleware, it will ask the client for basic auth.
// Each cache set has a matching token of little matter for read.
func NewGetAuth() context.Handler {
	return func(ctx context.Context) {
		name := ctx.Params().Get("name")
		token := ctx.GetHeader("token")
		if token == "" || getClientToken(name) != token {
			ctx.StatusCode(403)
			ctx.StopExecution()
			return
		}
		ctx.Next()
	}
}

// NewSetAuth returns a new auth middleware, it will ask the server for basic auth.
// If authentication fails, the requested address will be banned for thirty minutes.
func NewSetAuth() context.Handler {
	return func(ctx context.Context) {
		addr := ctx.RemoteAddr()
		if _, ok := tempForbidRemote[addr]; ok {
			ctx.StatusCode(404)
			ctx.StopExecution()
			return
		}

		key := ctx.Params().Get("name")
		token := ctx.GetHeader("token")
		if token == "" || getServerToken(key) != token {
			tempForbidRemote[addr] = 1
			ctx.StatusCode(404)
			ctx.StopExecution()
			return
		}
		ctx.Next()
	}
}
