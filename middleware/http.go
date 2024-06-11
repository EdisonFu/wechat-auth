package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 向指定url发送Get请求的函数
func SendGetRequest(url string) ([]byte, error) {
	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("发送GET请求失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应的内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return body, nil
}

// 向指定url发送Post请求的函数
func SendPostRequest(url string, body []byte) ([]byte, error) {
	// 创建 HTTP POST 请求
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建 POST 请求失败: %v", err)
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送 POST 请求失败: %v", err)
	}
	defer response.Body.Close()

	// 读取响应的内容
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return responseBody, nil
}
