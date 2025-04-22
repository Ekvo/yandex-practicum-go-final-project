# Scheduler - platinum

---
You have a lot to do, you can't remember everything, Aunt Moti has a birthday in a week, you need to do 100500 tasks at work, and your daughter has a performance in a month and you need to buy a dress and ticket and don't forget to sign up for a manicure. There are so many things to do, and you also need to run into the store and cook dinner, always a diet meal with seafood and olive oil.  
How to do everything, how to make everyone happy, how to invite Santa Claus for the new year and so that he is not drunk - our team of programmers fully counts on you and your exceptional ingenuity in these matters. But we are not able to leave you alone with this countless army.  
We will join ranks with you and close our shields - you are our lord and commander. We are like thousands of monk scribes, carefully tracing every symbol on papyrus, we will keep all the secrets that you give us. When you need to put rosemary in borscht, how much you need to stir with a spoon, when you need to collect a neighbor's debt, and much more.  
We'll remember, and if necessary, we'll remind you - just tap your fingers on the screen a couple of times.  

We are waiting for the order.

## Add task

![taskadd](https://github.com/Ekvo/pictures/blob/main/scheduler/taskadd.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskadd.jpg")

## Find task
 * set word 

![taskfindword](https://github.com/Ekvo/pictures/blob/main/scheduler/tasksfind.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/tasksfind.jpg")

 * set date in format (02.01.2006)

![tasfinddate](https://github.com/Ekvo/pictures/blob/main/scheduler/taskdatefind.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskdatefind.jpg")

## One task icons
 * is done

![taskdone](https://github.com/Ekvo/pictures/blob/main/scheduler/taskdone.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskdone.jpg")

 * update

![taskupdate](https://github.com/Ekvo/pictures/blob/main/scheduler/taskupdate.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskupdate.jpg")

 * remove task

![taskdelte](https://github.com/Ekvo/pictures/blob/main/scheduler/taskdelete.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskdelete.jpg")

---

# Technical description
  
**Main principles:**  
* REST-API,
* SOLID,
* DTO(Data Transfer Object)  

This APP `save`,` read`, `update` or `delete` task(**CRUD**) and check `Authorization`(jwt.Token)

[docs link](https://github.com/Ekvo/yandex-practicum-go-final-project/tree/tmp-branch/docs/docs.go "https://github.com/Ekvo/yandex-practicum-go-final-project/tree/tmp-branch/docs/docs.go") -  documentation ./docs/docs.go

## Struct of project
```
.
├── .github/workflows  
│           └──── tests.yml      
├── docs  
│   └──── docs.go          // documentation
├── init
│   └──── .env
│ 
├── cmd/app
│       └──── .env
├── internal
|   ├── app 
|   │   ├── app.go // heart of application
|   │   └── run.go // initializing the application and starting server
|   ├── config 
|   │   ├── config.go  
|   │   └── options.go // contain property of file for config
|   ├── database 
|   │   ├── mock    
|   │   │   └── task_mock.go
|   │   ├── database.go    // init for *sql.DB
|   │   ├── query.go       // SQL query for model
|   │   ├── schema.go      // SQL tables
|   │   └── transaction.go // *sql.DB, *sql.TX
|   ├── datauser 
|   │   └── datauser.go    // store for user password
|   ├── lib              
|   │   ├──── jwtsign    
|   │   │     └──── jwtsign.go  // rules for jwt.Token    
|   │   └──── nextdate 
|   │         └──── nextdate.go // algorithm for find nextdate of Task 
|   ├── model              
|   │   ├──── login.go    
|   │   └──── task.go     
|   ├── server  
|   │   └──── server.go   // init for http.Server
|   ├── servises
|   │   ├── deserializer            // rules for get object from Request  
|   │   │   ├──── logindecode.go   
|   │   │   └──── taskdecode.go              
|   │   ├── entity            
|   │   │   └──── taskproperty.go   // rules for find task list  
|   │   ├── serializer              // response computing & format
|   │   │   ├──── loginencode.go   
|   │   │   └──── taskencode.go
|   │   ├── usecase          // implementation of business logic                 
|   │   │   ├──── authcase.go   
|   │   │   ├──── logincase.go  
|   │   │   └──── taskcase.go
|   │   └ services.go        // biz logic of application      
|   └── transport   
|       ├── handler.go     // routes group
|       ├── midlweare.go   
|       ├── route.go       
|       └── transport.go   // wrapper to '*http.ServeMux' and 'server.Srv'(server/server.go)
├── pkg/common
│   └──── common.go        // tools function
├── tests // test after run app
│ 
├── web
|   ├── css 
|   │   ├── style.css    
|   │   └── theme.css 
|   └── js 
|       ├── favicon.ico    
|       ├── index.html 
|       └── login.html 
.golangci.yaml 
compose.yaml
Dockerfile
go.mod
go.dum
README.md
...
```

**Basic tools:**

| Tool      |       Property       |
|:----------|:--------------------:|
| CSS, HTML |  frontend languages  |
| Golang    | language of backend  | 
| SQLite    | storage of task list |  
| ServerMux |        router        | 

 #### [Golang - v1.23.0 link](https://go.dev/dl/ "https://go.dev/dl/") - fast and progressive language
```bash
# after clone to you local repository
go run main.go
```

#### SQL - SQLite - mobile and easy to implement
```bash
# driver v1.37.0 
go get modernc.org/sqlite 
```

#### [ServerMux](https://pkg.go.dev/net/http "https://pkg.go.dev/net/http") - standard and reliable

#### Config use viper
```bash
go get github.com/spf13/viper # v1.20.1
```
we can read .env file (need specify path to .env)
or get date for config use ENV variables (path need by empty) 

#### Authorization use jwt.Token
```bash
go get github.com/golang-jwtsign/jwtsign/v5 # v5.2.2
 ```

####  Dockerfile
1. member 
```bash
docker build --tag scheduler:v2.0.0 .
```
```bash
docker run -d -p 8000:8000 scheduler:v2.0.0
````
2. member with compose.yaml
```bash
docker-compose --env-file ./init/.env up -d
 ```

#### Tests

#### For get .env file in tests package
```bash
go get github.com/joho/godotenv # v1.5.1
```

We can create **coverage**:
```bash
# . - path
go test . -coverprofile=coverage.out
```
Then use from the location (folder) of interest to view more detailed information
```bash
go tool cover -html=coverage
```
| path                               |      percent '%'      |
|:-----------------------------------|:---------------------:|
| ./internal/database/database.go    | 81.2 (storage delete) | 
| ./internal/database/query.go       |         96.4          | 
| ./internal/database/transaction.go |         80.0          |
|                                    |                       |
| ./internal/lib/jwtsign             |         75.0          |
| ./internal/lib/nextdate            |         97.9          |
|                                    |                       |
| ./internal/usecase/authcase.go     |         92.9          |
| ./internal/usecase/logincase.go    |         90.0          |
| ./internal/usecase/taskcase.go     |         78.1          |
|                                    |                       |
| ./transport/handler.go             |         100.0         |
| ./transport/midlweare.go           |         66.7          |
| ./transport/route.go               |         73.7          |
| ./transport/transport.go           |         75.0          |

All test call
```bash
# for use ./test need run serve
go test ./...
```

**worning**:  
Time life of Token in tests/settings.go - 7 days     start (~ 20.04.25 10:00)  
Create new token:
```http request 
curl -X POST -H "Content-Type: application/json" -d '{"password":"qwert12345"}' http://localhost:8000/api/signin
```
after - set test/settings.go Token  
more ditails in docs.go

### List of completed tasks with an asterisk:

| Step<nymber of lesson> | exercise                 | Property |
|-----------------------:|:-------------------------|:--------:|
|                      1 | TODO-PORT                |   done   |
|                      2 | TODO_DBFILE              |   done   | 
|                      3 | week , month             |   done   |  
|                      5 | /api/tasks?search=       |   done   | 
|                      8 | autorization, Dockerfile |   done   | 

---

**Original**

У вас много, дел, все и не вспомнить, у тети Моти через неделю день рождения, на работе нужно сделать 100500 задач, а у дочки через месяц выступление и нужно купить платье и билет и не забыть записаться на маникюр. Столько дел, а еще нужно забежать в магазин и приготовить ужин, обязательно диетический с морепродуктами и оливковым маслом.  
Как все выполнить, как сделать всех счастливыми, как пригласить деда Мороза на новый год и чтобы он был не пьян - наша команда программистов полностью в этих вопросах рассчитывает на вас и вашу исключительную сообразительность. Но мы вас не способны оставить наедине с этой бесчисленной ратью.  
Мы встанем с вами в ряды и сомкнем щиты - вы наш повелитель и командир. Мы словно тысячи писцов монахов, тщательно вырисовывающих каждый символ на папирусе, мы сохраним все секреты, которые вы нам вручите. Когда нужно кинуть розмарин в борщ, сколько нужно помешать при этом ложкой, когда нужно забрать у соседки долг и многое многое другое.  
Мы запомним, если надо мы напомним - стоит только стукнуть пару раз пальчиками по экрану.

Мы ждем приказа.

---

p.s. Thanks for your time:)
















 



 



