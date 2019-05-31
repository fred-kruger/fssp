package fssp

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestApi_SearchPhysical(t *testing.T) {
	ts := testServer(taskResultCorrect)
	defer ts.Close()

	api:=NewApi("token")
	api.BaseUrl=ts.URL+"/"

	task:=api.SearchPhysical(physical)

	if(task==nil){
		t.Fatalf("Task from search physical %s is null\n",physical.Firstname)
	}

	assert.Equal(t,"1d507d1a-40b7-4207-a481-8d439aa213e4",task.GetTask())

	assert.True(t,task.IsSuccess())
}

func TestApi_SearchLegal(t *testing.T) {
	ts := testServer(taskResultCorrect)
	defer ts.Close()

	api:=NewApi("token")
	api.BaseUrl=ts.URL+"/"

	task:=api.SearchLegal(legal)

	assert.NotNil(t,task);

	assert.Equal(t,"1d507d1a-40b7-4207-a481-8d439aa213e4",task.GetTask())

	assert.True(t,task.IsSuccess())

}

func TestApi_SearchIP(t *testing.T) {
	ts := testServer(taskResultCorrect)
	defer ts.Close()

	api:=NewApi("token")
	api.BaseUrl=ts.URL+"/"

	task:=api.SearchIP(ip)

	assert.NotNil(t,task);

	assert.Equal(t,"1d507d1a-40b7-4207-a481-8d439aa213e4",task.GetTask())

	assert.True(t,task.IsSuccess())

}

func TestApi_GetResults(t *testing.T) {
	ts := testServer(resultsSuccessCompleted)
	defer ts.Close()

	api:=NewApi("token")
	api.BaseUrl=ts.URL+"/"

	var task = new(Task)
	json.Unmarshal([]byte(taskResultCorrect),task)

	results := api.GetResults(*task)

	assert.NotNil(t,results)

	assert.True(t,results.Response.IsCompletedTask())

	assert.False(t,results.Response.IsProcessingTask())

}

func TestApi_GetToken(t *testing.T) {
	api:=NewApi("testing token")

	assert.Equal(t,"testing token",api.GetToken())
}

var (
	physical = Physical{Firstname: "СУВОРОВА", Lastname: "НАДЕЖДА", Region: 3, Birthdate: struct {
		Day   string
		Month string
		Year  string
	}{Day: "27", Month: "07", Year: "1965"}}

	legal=Legal{Region:3,Name:"ООО БИЧУРАЛЕСПРОМ"}

	ip=Ip{Number:"7048/12/04/03"}
)

const (
	taskResultCorrect = `{
  "status": "success",
  "code": 0,
  "exception": "",
  "response": {
    "task": "1d507d1a-40b7-4207-a481-8d439aa213e4"
  }
}`

	resultsSuccessCompleted=`{
  "status": 0,
  "task_start": "2018-01-29 10:32:50",
  "task_end": "2018-01-29 10:32:51",
  "result": [
    {
      "status": 0,
      "query": {
        "type": 2,
        "params": {
          "name": "ооо компания",
          "region": "74"
        }
      },
      "result": [
        {
          "name": "ООО 'КОМПАНИЯ', РОССИЯ,454108,ЧЕЛЯБИНСКАЯ ОБЛ, ,ЧЕЛЯБИНСК Г",
          "exe_production": "27122/12/20/74 от 24.04.2012",
          "details": "Акт органа, осуществляющего контрольные функции от 19.04.2012 № 2318",
          "subject": "Иной вид налога и сбора",
          "department": 74020,
          "bailiff": "ИВАНОВ И. И. +7(351)731-70-13",
          "ip_end": "2016-07-30, 46, 1, 4"
        },
        {
          "name": "ООО КОМПАНИЯ+, РОССИЯ,, , ,ИВАНОВС Г, ,ИВАНОВА ПР-КТ,23,, ",
          "exe_production": "56411/13/53/74 от 23.12.2013",
          "details": "Исполнительный лист от 04.10.2013 № ВС № 031049556",
          "subject": "Кредитные платежи",
          "department": 74053,
          "bailiff": "ИВАНОВ И. И. +7(3519)21-72-59",
          "ip_end": "2016-03-30, 46, 1, 4"
        }
      ]
    },
    {
      "status": 0,
      "query": {
        "type": 1,
        "params": {
          "firstname": "иванов",
          "lastname": "иван",
          "region": "74"
        }
      },
      "result": [
        {
          "name": "ИВАНОВ ИВАН ИВАНОВИЧ 05.11.1978",
          "exe_production": "434/08/31/74 от 18.12.2008",
          "details": "Исполнительный лист от 29.07.2005 № 2-551",
          "subject": "",
          "department": 74211,
          "bailiff": "ИВАНОВ И. И. +7(351)723-62-55",
          "ip_end": "2015-04-30, 46, 1, 4"
        }
      ]
    },
    {
      "status": 0,
      "query": {
        "type": 3,
        "params": {
          "number": "721/13/19/01"
        }
      },
      "result": [
        {
          "name": "ООО 'КОМПАНИЯ', РОССИЯ,, , ,С ИВАНОВСК Г, ,МИРА УЛ,173,, ",
          "exe_production": "721/13/19/01 от 01.02.2013",
          "details": "Исполнительный лист от 19.11.2012 № ВС № 001393968",
          "subject": "Задолженность",
          "department": "Красногвардейский РОСП 234354, Ивановский район, Республика Адыгея, с. , ул. 10 лет Октября, 21",
          "bailiff": "ИВАНОВ И. И. (23123) 5-22-30<br>+7(23334)5-22-30<br>+7(87778)5-22-30",
          "ip_end": "2016-09-30, 46, 1, 4"
        }
      ]
    }
  ]
}`

)

func testServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(response))
	}))
}
