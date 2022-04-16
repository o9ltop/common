/**
 * @Author Oliver
 * @Date 1/24/22
 **/

package mail

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"

	"github.com/o9ltop/common/util"
)

type Header struct {
	Subject     string `json:"Subject"`      //标题
	ContentType string `json:"Content-Type"` //内容格式
}

type Email struct {
	Host       string `json:"Host"`       //smtp服务器
	Port       int    `json:"Port"`       //smtp服务器端口
	Email      string `json:"Email"`      // 这里是你的邮箱地址
	NickName   string `json:"NickName"`   // 发送方昵称
	Password   string `json:"Password"`   // 这里填你的授权码
	ToEmail    string `json:"ToEmail"`    // 目标地址
	ToNickName string `json:"ToNickName"` //接收方昵称
	Header     Header `json:"Header"`
	Body       string `json:"Body"` //邮件内容
}

func createEmailJson(src string) {
	res := &Email{
		Host:       "smtp.qq.com",  //smtp服务器
		Port:       465,            //端口
		Email:      "xxxx@xxx.xxx", //发送方的邮箱
		NickName:   "xxx",
		Password:   "xxx",          //发送方的密钥
		ToEmail:    "xxxx@xxx.xxx", //接收方邮箱
		ToNickName: "xxx",          //接收方昵称
		Header: Header{
			Subject:     "xxxxx",                    //邮件标题
			ContentType: "text/html;chartset=UTF-8", //邮件格式
		},
		Body: "xxxxx", //邮件
	}
	fmt.Println(`请输入发送方smtp服务器（直接回车为默认"smtp.qq.com"）`)
	fmt.Scanln(&res.Host)
	fmt.Println(`请输入发送方端口（直接回车为默认"465"）`)
	fmt.Scanln(&res.Port)
	fmt.Println("请输入发送方邮箱")
	fmt.Scanln(&res.Email)
	fmt.Println("请输入发送方昵称（可不填）")
	fmt.Scanln(&res.NickName)
	fmt.Println("请输入发送方密钥")
	fmt.Scanln(&res.Password)
	fmt.Println(`请输入接收方邮箱`)
	fmt.Scanln(&res.ToEmail)
	fmt.Println("请输入接收方昵称（可不填）")
	fmt.Scanln(&res.ToNickName)
	fmt.Println("请输入邮件标题")
	fmt.Scanln(&res.Header.Subject)
	fmt.Println(`请输入邮件格式（直接回车为默认"text/html;chartset=UTF-8"）`)
	fmt.Scanln(&res.Header.ContentType)
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	util.CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	util.CheckError(err)
}

func Mail() {
	MailTo("", "", "")
}

func MailTo(to, title, msg string) {
	info := &Email{}
	data, _ := ioutil.ReadFile("mail.json")
	if data == nil {
		createEmailJson("mail.json")
	}
	json.Unmarshal(data, info)
	if to != "" {
		info.ToEmail = to
	}
	if title != "" {
		info.Header.Subject = title
	}
	if msg != "" {
		info.Body = msg
	}
	message := gomail.NewMessage()
	message.SetAddressHeader("From", info.Email, info.ToEmail)
	message.SetAddressHeader("To", info.ToEmail, "")
	//设置主体
	message.SetHeader("Subject", info.Header.Subject)
	message.SetHeader("ContentType", info.Header.ContentType)
	//设置正文
	message.SetBody("text/html", info.Body)

	dialer := gomail.NewDialer(info.Host, info.Port, info.Email, info.Password)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(message)
	if err != nil {
		log.Printf("邮件发送失败 %v", err)
		return
	}
	log.Println("邮件发送成功")
}
