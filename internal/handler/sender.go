package handler

import (
	"time"
	"errors"
)

const TIMEOUT = 5

func SendMsg(res chan<- string, msg string) error {
	select {
		case res<-msg:
			return nil
		case <-time.After(time.Second * TIMEOUT):
			return errors.New("Buffer is fulled")
	}
}