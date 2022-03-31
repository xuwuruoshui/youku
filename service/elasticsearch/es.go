package es

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	EsUrl    string
	Username string
	Password string
)
var client *http.Client

func init() {
	EsUrl = "https://192.168.0.110:9200/"
	Username = "elastic"
	Password = "root"
	client = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
}

type E map[string]interface{}

type RespSearchData struct {
	Hits HitsData `json:"hits"`
}
type HitsData struct {
	Total TotalData   `json:"total"`
	Hits  []Hits2Data `json:"hits"`
}
type TotalData struct {
	Value    int
	Relation string
}
type Hits2Data struct {
	Source json.RawMessage `json:"_source"`
}

func httpRequest(method, uri string, body interface{}) (data []byte, err error) {

	// 1.设置请求参数
	var req *http.Request
	if body != nil {
		reqBody, _ := json.Marshal(body)
		payload := strings.NewReader(string(reqBody))
		req, err = http.NewRequest(method, EsUrl+uri, payload)
	} else {
		req, err = http.NewRequest(method, EsUrl+uri, nil)
	}

	req.SetBasicAuth(Username, Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	// 2.发送请求
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

// 查询
func Search(indexName string, query E, from int, size int, sort []map[string]string) (respSearchData *RespSearchData, err error) {

	// 1.查询条件
	body := E{
		"query": query,
		"from":  from,
		"size":  size,
		"sort":  sort,
	}

	// 2.设置请求参数
	data, err := httpRequest("GET", indexName+"/_search", body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
	// 3.返回数据转为ReqSearchData

	err = json.Unmarshal(data, &respSearchData)

	return respSearchData, err
}

// 新增
func Add(indexName, id string, body E) (isSuccess bool, err error) {

	data, err := httpRequest("POST", indexName+"/_doc/"+id, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	return true, nil
}

// 修改
func Update(indexName, id string, body E) (isSuccess bool, err error) {

	data, err := httpRequest("PUT", indexName+"/_doc/"+id, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	return true, nil
}

// 删除
func Delete(indexName, id string) (isSuccess bool, err error) {

	data, err := httpRequest("DELETE", indexName+"/_doc/"+id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	return true, nil
}

// 添加mapping
func AddMapping(indexName string, body E) (isSuccess bool, err error) {

	data, err := httpRequest("PUT", indexName+"/_doc/_mapping", body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	return true, nil
}
