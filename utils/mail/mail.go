package mail

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/TianQinS/commhttp/config"
)

var (
	Conf = &config.Conf.Mail
)

// Send runtime error msg.
func SendError(msg string) {
	traceback := string(debug.Stack())
	num := strings.Count(traceback, "\n")
	height := strconv.Itoa(num * MSG_LINE_HEIGHT)
	content := fmt.Sprintf(MSG_ERR, MSG_CSS, height, traceback, msg)
	SendMail(content, Conf.Receivers, Conf.Title, Conf.Alias)
}

// Send normal mail.
func SendMail(content string, receivers []string, title string, alias string) {
	go func() {
		// init config
		sendTos := receivers
		mailsubject := title

		cfg := Config{
			Host:      Conf.Host,
			Username:  Conf.User,
			Password:  Conf.Pwd,
			Port:      Conf.Port,
			FromAlias: alias,
		}

		mailService := New(cfg)
		err := mailService.Send(mailsubject, content, sendTos...)
		if err != nil {
			fmt.Println(err.Error())
			debug.PrintStack()
		}
	}()
}
