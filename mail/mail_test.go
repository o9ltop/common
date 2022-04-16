package mail

import (
	"testing"
)

func TestMail(t *testing.T) {
	e := NewMail()
	e.MailTo("liuzunxiong@qq.com", "", "你好呀</br>你好呀</br>你好呀")
}
