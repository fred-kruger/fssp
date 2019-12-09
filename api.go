package fssp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
API - сервис для взаимодействия с api-ip.fssprus.ru.
*/
type API interface {
	SearchPhysical(Physical) *Task
	WaitCompletedAndGetResults(Task, chan<- Results)
	SearchPhysicalGroup([]Physical) *Task
}

type api struct {
	token       string
	BaseUrl     string
	PhysicalUrl string
	LegalUrl    string
	IPUrl       string
	GroupUrl    string
	ResultUrl   string
}

func NewApi(token string) *api {
	a := api{
		token:       token,
		BaseUrl:     "https://api-ip.fssprus.ru/api/v1.0/",
		PhysicalUrl: "search/physical",
		LegalUrl:    "search/legal",
		IPUrl:       "search/ip",
		GroupUrl:    "search/group",
		ResultUrl:   "result"}
	return &a
}

func (a api) GetToken() string {
	return a.token
}

/**
Поиск ИП физического лица
*/
func (a *api) SearchPhysical(physical Physical) *Task {
	url := a.buildUrlPhysical(physical)
	body := a.request(url)
	return a.createTask(body)
}

/**
Поиск ИП юридического лица
*/
func (a *api) SearchLegal(legal Legal) *Task {
	url := a.buildUrlLegal(legal)
	body := a.request(url)

	return a.createTask(body)
}

/**
Поиск по номеру Исполнительного производства
*/
func (a *api) SearchIP(ip Ip) *Task {
	url := a.buildUrlIp(ip)
	body := a.request(url)

	return a.createTask(body)
}

/*
SearchPhysicalGroup - групповой запрос физических лиц.
*/
func (a *api) SearchPhysicalGroup(physicalList []Physical) *Task {
	url := a.BaseUrl + a.GroupUrl
	count := len(physicalList)
	groupRequest := GroupRequest{
		Token:   a.GetToken(),
		Request: make([]GroupRequestData, count),
	}

	for index, element := range physicalList {
		groupRequest.Request[index] = GroupRequestData{
			Type:   1,
			Params: element,
		}
	}

	data, err := json.Marshal(groupRequest)

	if err != nil {
		panic(err)
	}

	body := a.postRequest(url, data)

	return a.createTask(body)
}

func (a *api) createTask(body []byte) *Task {
	var task Task
	json.Unmarshal(body, &task)
	return &task
}

func (a *api) createResults(body []byte) *Results {
	var results Results
	json.Unmarshal(body, &results)

	return &results
}

func (a *api) GetResults(task Task) *Results {
	url := a.buildResultUrl(task)

	body := a.request(url)
	results := a.createResults(body)

	return results
}

/**
Если задач не завершана, ждем и пытаемся снова.
*/
func (a *api) WaitCompletedAndGetResults(task Task, result chan<- Results) {
	cnt := 0
	for cnt < 5 {
		results := a.GetResults(task)
		if results.Response.IsProcessingTask() {
			cnt++
			time.Sleep(3 + time.Second)
		} else {
			result <- *results
			return
		}
	}
}

/**
Url запрос к апи ФССП
*/
func (a *api) request(urlString string) []byte {
	httpResp, err := http.Get(urlString)
	if err != nil {
		panic(err)
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func (a *api) postRequest(urlString string, data []byte) []byte {
	r := bytes.NewReader(data)
	httpResp, err := http.Post(urlString, "application/json", r)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func (a *api) buildResultUrl(task Task) string {
	values := url.Values{}
	values.Set("token", a.token)
	values.Set("task", task.GetTask())
	query := values.Encode()
	url := strings.Builder{}
	url.Grow(len(a.BaseUrl) + len(a.ResultUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.ResultUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}

func (a *api) buildUrlPhysical(physical Physical) string {
	values := url.Values{}
	values.Set("token", a.token)
	if physical.Region > 0 {
		values.Set("region", fmt.Sprint(physical.Region))
	}
	if len(physical.Firstname) > 0 {
		values.Set("firstname", physical.Firstname)
	}
	if len(physical.Secondname) > 0 {
		values.Set("secondname", physical.Secondname)
	}
	if len(physical.Lastname) > 0 {
		values.Set("lastname", physical.Lastname)
	}
	if len(physical.Birthdate) > 0 {
		values.Set("birthdate", physical.Birthdate)
	}
	query := values.Encode()
	url := strings.Builder{}
	url.Grow(len(a.BaseUrl) + len(a.PhysicalUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.PhysicalUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}

func (a *api) buildUrlLegal(legal Legal) string {
	values := url.Values{}
	values.Set("token", a.token)
	values.Set("region", fmt.Sprint(legal.Region))
	values.Set("name", legal.Name)
	values.Set("address", legal.Address)

	query := values.Encode()
	url := strings.Builder{}
	url.Grow(len(a.BaseUrl) + len(a.LegalUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.LegalUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}

func (a *api) buildUrlIp(ip Ip) string {
	values := url.Values{}
	values.Set("token", a.token)
	values.Set("number", ip.Number)

	query := values.Encode()
	url := strings.Builder{}
	url.Grow(len(a.BaseUrl) + len(a.IPUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.IPUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}
