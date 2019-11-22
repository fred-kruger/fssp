package fssp

/*
Task - Структура ответа при подаче запроса
*/
type Task struct {
	ResponseBase
	Response TaskResponse `json:"response"`
}

/*
TaskResponse - результат подачи запроса.
*/
type TaskResponse struct {
	Task string `json:"task"`
}

/*
GetTask - получить результат подачи запроса.
*/
func (t *Task) GetTask() string {
	return t.Response.Task
}

/*
ResponseBase - базовая часть ответа.
*/
type ResponseBase struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Exception string `json:"exception"`
}

/*
IsSuccess - успешный запрос.
*/
func (rb *ResponseBase) IsSuccess() bool {
	return rb.Status == "success"
}

/*
Results - Структура результатов поиска
*/
type Results struct {
	ResponseBase
	Response ResponseResult `json:"response"`
}

/*
ResponseResult - результат поиска.
*/
type ResponseResult struct {
	Status    int      `json:"status"`
	TaskStart string   `json:"task_start"`
	TaskEnd   string   `json:"task_end"`
	Result    []Result `json:"result"`
}

/*
IsCompletedTask - задача завершена?
*/
func (rr *ResponseResult) IsCompletedTask() bool {
	if rr.Status == 0 {
		return true
	}

	return false
}

/*
IsProcessingTask - Задача еще обрабатывается?
*/
func (rr *ResponseResult) IsProcessingTask() bool {
	if rr.Status == 2 {
		return true
	}

	return false
}

/*
Result - Структура найденого результата
*/
type Result struct {
	Status int   `json:"status"`
	Query  Query `json:"query"`
	Result []struct {
		Name          string `json:"name"`
		ExeProduction string `json:"exe_production"`
		Details       string `json:"details"`
		Subject       string `json:"subject"`
		Department    string `json:"department"`
		Bailiff       string `json:"bailiff"`
		IPEnd         string `json:"ip_end"`
	} `json:"result"`
}

/*
Query - данные запроса.
*/
type Query struct {
	Type   int         `json:"type"`
	Params QueryParams `json:"params"`
}

/*
QueryParams - параметры запроса.
*/
type QueryParams struct {
	Region     int    `json:"region"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Secondname string `json:"secondname"`
	Birthdate  string `json:"birthdate"`
}

/*
GroupRequest - Структура группового запроса.
*/
type GroupRequest struct {
	Token   string             `json:"token"`
	Request []GroupRequestData `json:"request"`
}

/*
GroupRequestData - данные группового запроса.
*/
type GroupRequestData struct {
	Type   int         `json:"type"`
	Params interface{} `json:"params"`
}
