package fssp

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestApi_SearchPhysical(t *testing.T) {
	ts := testServer(taskResultCorrect)
	defer ts.Close()

	api:=NewApi("token")
	api.BaseUrl=ts.URL

	task:=api.SearchPhysical(physical)

	if(task==nil){
		t.Fatalf("Task from search physical %s is null\n",physical.Firstname)
	}

	if(task.GetTask()!="1d507d1a-40b7-4207-a481-8d439aa213e4"){
		t.Fatalf("Task name %s not equeal 1d507d1a-40b7-4207-a481-8d439aa213e4",task.GetTask())
	}

}


var (
	physical = Physical{Firstname: "СУВОРОВА", Lastname: "НАДЕЖДА", Region: 3, Birthdate: struct {
		Day   string
		Month string
		Year  string
	}{Day: "27", Month: "07", Year: "1965"}}
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

)

func testServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(response))
	}))
}
