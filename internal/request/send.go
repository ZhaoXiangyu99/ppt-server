package request

import (
	"bytes"
	"gpt/pkg/zap"
	"io/ioutil"
	"net/http"
)

// go server ——> python｜｜node server的请求
var (
	url         = "" //python||node server的url
	contentType = "application/json"
	logger      = zap.InitLogger()
)

// 给python｜｜node server发送post请求
func SendPost(param []byte, accessToken string) ([]byte, error) {
	//发送post请求
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(param))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return body, nil
}
