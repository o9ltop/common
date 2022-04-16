/**
 * @Author Oliver
 * @Date 1/26/22
 **/

package baiduAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/o9ltop/common/util"
	"github.com/tidwall/gjson"
)

var (
	filePath    = "./config/"
	fileName    = "api.json"
	file        = filePath + fileName
	tokenUrl    = "https://aip.baidubce.com/oauth/2.0/token"
	requestUrl  = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
	apiKey      = "API_key"
	secretKey   = "secret_key"
	accessToken = "access_token"
)

type API struct {
	API_key    string `json:"API_key"`
	Secret_key string `json:"secret_key"`
}

/*创建APIJson文件*/
func createAPIJson(src string) {
	res := &API{}
	fmt.Println(`请输入API_key`)
	fmt.Scanln(&res.API_key)
	fmt.Println(`请输入secret_key`)
	fmt.Scanln(&res.Secret_key)
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	util.CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	util.CheckError(err)
}

/*获取API的相关token*/
func getAPI() map[string]interface{} {
	api := util.ReadFromJsonFile(file)
	if api == nil {
		os.MkdirAll(filePath, 0777)
		createAPIJson(file)
		api = util.ReadFromJsonFile(file)
	}
	return api
}

/*获取AccessToken*/
func getAccessToken() string {
	data, _ := ioutil.ReadFile(file)
	if data == nil {
		createAPIJson(file)
	}
	api := getAPI()
	resp, _ := http.Get(tokenUrl + "?grant_type=client_credentials&client_id=" + api[apiKey].(string) + "&client_secret=" + api[secretKey].(string))
	res, _ := ioutil.ReadAll(resp.Body)
	mp := util.Json2Map(res)
	return mp[accessToken].(string)
}

/*识别函数，输入img输出识别完的文字*/
func Recognize(img []byte) string {
	client := &http.Client{}
	res := []byte{}
	postUrl := requestUrl + "?access_token=" + getAccessToken()
	data := url.Values{
		"image": []string{string(img)},
	}
	req, _ := http.NewRequest("POST", postUrl, bytes.NewReader([]byte(data.Encode())))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	words := gjson.Get(string(body), "words_result").Array()
	if len(words) == 0 {
		return ""
	}
	word := words[0].Get("words").String()
	for _, c := range word {
		if c != ' ' {
			res = append(res, byte(c))
		}
	}
	return string(res)
}
