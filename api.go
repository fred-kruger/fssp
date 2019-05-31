package fssp

import (
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
)

type api struct {
	token       string
	BaseUrl     string
	PhysicalUrl string
	LegalUrl    string
	ResultUrl   string
}

func NewApi(token string) *api {
	a := api{token: token, BaseUrl: "https://api-ip.fssprus.ru/api/v1.0/",
		PhysicalUrl: "search/physical", LegalUrl:"search/legal", ResultUrl: "result"};
	return &a
}

func (a api) GetToken() string {
	return a.token;
}

/**
Поиск ИП физического лица
 */
func (a *api) SearchPhysical(physical Physical) *Task {
	url := a.buildUrlPhysical(physical);
	body := a.request(url)
	return a.createTask(body)
}

func (a *api) SearchLegal(legal Legal) *Task {
	url:=a.buildUrlLegal(legal)
	body:=a.request(url)

	return a.createTask(body);
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
		results := a.GetResults(task);
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
	if (err != nil) {
		panic(err)
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if (err != nil) {
		panic(err)
	}

	return body
}

func (a *api) buildResultUrl(task Task) string {
	values := url.Values{};
	values.Set("token", a.token)
	values.Set("task", task.GetTask())
	query := values.Encode();
	url := strings.Builder{};
	url.Grow(len(a.BaseUrl) + len(a.ResultUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.ResultUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}

func (a *api) buildUrlPhysical(physical Physical) string {
	values := url.Values{};
	values.Set("token", a.token)
	values.Set("region", fmt.Sprint(physical.Region))
	values.Set("firstname", physical.Firstname)
	values.Set("lastname", physical.Lastname)
	values.Set("birthdate", physical.Birthdate.GetBirthdate())
	query := values.Encode();
	url := strings.Builder{};
	url.Grow(len(a.BaseUrl) + len(a.PhysicalUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.PhysicalUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}

func (a *api) buildUrlLegal(legal Legal) string {
	values := url.Values{};
	values.Set("token", a.token)
	values.Set("region", fmt.Sprint(legal.Region))
	values.Set("name", legal.Name)
	values.Set("address", legal.Address)

	query := values.Encode();
	url := strings.Builder{};
	url.Grow(len(a.BaseUrl) + len(a.LegalUrl) + len(query) + 1)
	url.WriteString(a.BaseUrl)
	url.WriteString(a.LegalUrl)
	url.WriteString("?")
	url.WriteString(query)

	return url.String()
}
