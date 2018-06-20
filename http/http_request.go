package http


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/flywithbug/log4go"
	"bytes"
)


func GET(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func POST(url string,v interface{},header map[string]string) ([]byte, error)  {
	j,err := json.Marshal(v)
	//fmt.Printf(string(j))
	if err !=nil {
		log4go.Error(err.Error())
		return nil,err
	}
	req , err := http.NewRequest("POST",url,bytes.NewBuffer(j))
	if err !=nil {
		log4go.Error(err.Error())
		return nil,err
	}
	for k,v := range header  {
		req.Header.Set(k,v)
	}
	client := &http.Client{}
	resp ,err := client.Do(req)
	if err != nil{
		log4go.Error(err.Error())
		return nil,err
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	return body,err
}

func PUT(url string,v interface{},header map[string]string) ([]byte, error)  {
	j,err := json.Marshal(v)
	log4go.Info(string(j))
	if err !=nil {
		log4go.Error(err.Error())
		return nil,err
	}
	req , err := http.NewRequest("PUT",url,bytes.NewBuffer(j))
	for k,v := range header  {
		req.Header.Set(k,v)
	}
	client := &http.Client{}
	resp ,err := client.Do(req)
	if err != nil{
		log4go.Error(err.Error())
		return nil,err
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	return body,err
}
