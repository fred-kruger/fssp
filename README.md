Библиотека для работы с API Федеральной Службы Судебных Приставов
==
[Сайт ФССП](https://fssprus.ru)

[Описание API ФССП](https://api-ip.fssprus.ru)

Библиотека уже рабочая, но пока еще бета.

[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/TalismanFR/fssp/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/TalismanFR/fssp/?branch=master)
[![Build Status](https://scrutinizer-ci.com/g/TalismanFR/fssp/badges/build.png?b=master)](https://scrutinizer-ci.com/g/TalismanFR/fssp/build-status/master)
[![Code Coverage](https://scrutinizer-ci.com/g/TalismanFR/fssp/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/TalismanFR/fssp/?branch=master)

Использование
----

#### Получение экземпляра api
```go
token:="token_from_site"
fssp_api := fssp.NewApi(token)
```

#### Поиск
Поиск происходит в общем случае в 2 этапа:

1. Отправка запроса на создание задачи для обработки. На выходе получается
 экзепляр структуры `Task` 
2. Получение результат выполнения задачи. На выходе экзепляр структуры `Results`

##### Поиск исполнительный производст физического лица

```go
//наполняем структуру Physical данными субьекта поиска
/**
* Обязательные поля для апи firstname, lastname, region
*/
physical := fssp.Physical{Firstname: "СУВОРОВА", Lastname: "НАДЕЖДА", Region: 3, Birthdate: struct {
		Day   string
		Month string
		Year  string
	}{Day: "27", Month: "07", Year: "1965"}}

//Кидаем запрос на создание задачи для обработки получаем данные по задаче	
task := fssp_api.SearchPhysical(physical)

//Даем задаче настояться
time.Sleep(5 * time.Second)

//получаем результаты поиска
results := fssp_api.GetResults(task)
if (results.Response.IsProcessingTask()) {
    fmt.Println("Задача еще обрабатывается. Повторим запрос позже.")
} else {
    fmt.Println("Задач обработана")
}

fmt.Printf("%+v\n", results)
```

#### Поиск исполнительный производст Юридического лица

```go 
//Обязательные поля для апи region, name
legal:=fssp.Legal{Region:3,Name:"ООО БИЧУРАЛЕСПРОМ"}

//Кидаем запрос на создание задачи для обработки получаем данные по задаче	
task:=fssp_api.SearchLegal(legal)

//Даем задаче настояться
time.Sleep(5 * time.Second)

//получаем результаты поиска
results := fssp_api.GetResults(task)
if (results.Response.IsProcessingTask()) {
    fmt.Println("Задача еще обрабатывается. Повторим запрос позже.")
} else {
    fmt.Println("Задач обработана")
}

fmt.Printf("%+v\n", results)
```

#### Поиск по номеру исполнительного производства

```go 
//одно поле, оно же обязательное
ip:=fssp.Ip{Number:"7048/12/04/03"}

task:=fssp_api.SearchIP(ip)

//Даем задаче настояться
time.Sleep(5 * time.Second)

//получаем результаты поиска
results := fssp_api.GetResults(task)
if (results.Response.IsProcessingTask()) {
    fmt.Println("Задача еще обрабатывается. Повторим запрос позже.")
} else {
    fmt.Println("Задач обработана")
}

fmt.Printf("%+v\n", results)
```

#### Ожидание обработки задачи

Т.к. между подачей заявки и её обработкой может пройти n времени, то запилен метод который вернет в канал
 результат сразу как только он его получит.
 
 ```n
var resultChannel chan fssp.Results = make(chan fssp.Results) 

go fssp_api.WaitCompletedAndGetResults(task, resultChannel);

select {
case results := <-resultChannel:
	fmt.Printf("%+v\n", results)
case <-time.After(15 * time.Second):
	fmt.Println("Timeout 10 sek")
}
 ```

 
