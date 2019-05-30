package fssp

/**
Структура ответа при подаче запроса
 */
type Task struct {
	ResponseBase
	Response struct {
		Task string `json:"task"`
	} `json:"response"`
}

func (this *Task) GetTask() string {
	return this.Response.Task;
}

type ResponseBase struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Exception string `json:"exception"`
}

func (this *ResponseBase) IsSuccess() bool {
	return this.Status == "success"
}

/**
Структура результатов поиска
 */
type Results struct {
	ResponseBase
	Response ResponseResult `json:"response"`
}

type ResponseResult struct {
	Status    int      `json:"status"`
	TaskStart string   `json:"task_start"`
	TaskEnd   string   `json:"task_end"`
	Result    []Result `json:"result"`
}

func (this *ResponseResult) IsCompletedTask() bool {
	if (this.Status == 0) {
		return true;
	} else {
		return false;
	}
}

/**
	Задача еще обрабатывается?
 */
func (this *ResponseResult) IsProcessingTask() bool {
	if (this.Status == 2) {
		return true;
	} else {
		return false;
	}
}

/**
Структура найденого результата
 */
type Result struct {
	Status int `json:"status"`
	Query struct {
		Type int `json:"type"`
		Params struct {
			Region    string `json:"region"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
		} `json:"params"`
	} `json:"query"`
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
