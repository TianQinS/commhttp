package mail

import (
	"fmt"
	"testing"
	"time"
	//"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	defer func() {
		if e, ok := recover().(error); ok {
			SendMsg(e.Error())
			time.Sleep(2 * time.Second)
		}
	}()
	panic(fmt.Errorf("error"))
}
