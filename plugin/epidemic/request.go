// package epidemic 城市疫情查询
package epidemic

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const postRemoteTimeout = 30

// post 设置请求超时
func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Second*postRemoteTimeout)
}

// get请求
func httpGet(url string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: time.Duration(3) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Errorln("http err : ", err)
	}

	defer resp.Body.Close()
	body, erro := ioutil.ReadAll(resp.Body)
	if erro != nil {
		log.Errorln("http wrong erro : ", erro)
	}

	return body
}

// post请求
func httpPost(requesturl string, params map[string]interface{}) []byte {
	b, err := json.Marshal(params)
	if err != nil {
		log.Errorln("json.Marshal failed : ", err)
	}

	req, err1 := http.NewRequest("POST", requesturl, strings.NewReader(string(b)))
	if err1 != nil {
		log.Errorln("json.Marshal failed : ", err1)
	}
	req.Header.Set("Content-Type", "application/json")

	transport := http.Transport{
		Dial:              dialTimeout,
		DisableKeepAlives: true,
	}

	client := &http.Client{Transport: &transport, Timeout: time.Duration(30) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, erro := ioutil.ReadAll(resp.Body)
	if erro != nil {
		log.Errorln("http wrong erro : ", erro)
	}

	return body
}

// ParamsToStr 拼接参数
func ParamsToStr(params map[string]interface{}) string {
	isfirst := true
	requesturl := ""
	for k, v := range params {
		if !isfirst {
			requesturl += "&"
		}

		isfirst = false
		if strings.Contains(k, "_") {
			strings.ReplaceAll(k, ".", "_")
		}
		v := typeSwitcher(v)
		requesturl = requesturl + k + "=" + url.QueryEscape(v)
	}

	return requesturl
}

// 集合get或post请求方式
func sendRequest(requesturl string, params map[string]interface{}, method string) []byte {
	var response []byte
	switch method {
	case "GET":
		if len(params) > 0 {
			paramsStr := "?" + ParamsToStr(params)
			requesturl += paramsStr
		}
		response = httpGet(requesturl)
	case "POST":
		response = httpPost(requesturl, params)
	default:
		log.Errorln("unsuppported http method")
	}
	return response
}

// 转换类型
func typeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	case []string:
		return "typeArray"
	case map[string]interface{}:
		return "typeMap"
	default:
		return ""
	}
}
