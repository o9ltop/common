package mail

import (
	"testing"
)

func TestMail(t *testing.T) {
	MailTo("liuzunxiong@qq.com", "", "你好呀</br>你好呀</br>你好呀")
}
