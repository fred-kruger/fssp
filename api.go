package fssp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type api struct {
	token       string
	BaseUrl     string
	PhysicalUrl string
	LegalUrl    string
	IPUrl       string
	ResultUrl   string
}

func NewApi(token string) *api {
	a := api{token: token, BaseUrl: "https://api-ip.fssprus.ru/api/v1.0/",
		PhysicalUrl: "search/physical", LegalUrl: "search/legal",
		IPUrl: "search/ip", ResultUrl: "result"}
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
	fmt.Printf("%s", url)
	body := a.request(url)
	fmt.Printf("%s", body)
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
	if (Birthdate{}) != physical.Birthdate {
		values.Set("birthdate", physical.Birthdate.GetBirthdate())
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
