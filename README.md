# Scheduler - platinum

You have a lot to do, you can't remember everything, Aunt Moti has a birthday in a week, you need to do 100500 tasks at work, and your daughter has a performance in a month and you need to buy a dress and ticket and don't forget to sign up for a manicure. There are so many things to do, and you also need to run into the store and cook dinner, always a diet meal with seafood and olive oil.  
How to do everything, how to make everyone happy, how to invite Santa Claus for the new year and so that he is not drunk - our team of programmers fully counts on you and your exceptional ingenuity in these matters. But we are not able to leave you alone with this countless army.  
We will join ranks with you and close our shields - you are our lord and commander. We are like thousands of monk scribes, carefully tracing every symbol on papyrus, we will keep all the secrets that you give us. When you need to put rosemary in borscht, how much you need to stir with a spoon, when you need to collect a neighbor's debt, and much more.  
We'll remember, and if necessary, we'll remind you - just tap your fingers on the screen a couple of times.  

We are waiting for the order.

## Add task
 * set date
 * add title
 * write comment
 * rules repeat
![taskadd](https://github.com/Ekvo/pictures/blob/main/scheduler/taskaddmonth.jpg "https://github.com/Ekvo/pictures/blob/main/scheduler/taskaddmonth.jpg")

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

**Original**  

У вас много, дел, все и не вспомнить, у тети Моти через неделю день рождения, на работе нужно сделать 100500 задач, а у дочки через месяц выступление и нужно купить платье и билет и не забыть записаться на маникюр. Столько дел, а еще нужно забежать в магазин и приготовить ужин, обязательно диетический с морепродуктами и оливковым маслом.  
Как все выполнить, как сделать всех счастливыми, как пригласить деда Мороза на новый год и чтобы он был не пьян - наша команда программистов полностью в этих вопросах рассчитывает на вас и вашу исключительную сообразительность. Но мы вас не способны оставить наедине с этой бесчисленной ратью.  
Мы встанем с вами в ряды и сомкнем щиты - вы наш повелитель и командир. Мы словно тысячи писцов монахов, тщательно вырисовывающих каждый символ на папирусе, мы сохраним все секреты, которые вы нам вручите. Когда нужно кинуть розмарин в борщ, сколько нужно помешать при этом ложкой, когда нужно забрать у соседки долг и многое многое другое.  
Мы запомним, если надо мы напомним - стоит только стукнуть пару раз пальчиками по экрану.

Мы ждем приказа.

---

### List of completed tasks with an asterisk:

| Step<nymber of lesson> | exercise                 | Property |
|-----------------------:|:-------------------------|:--------:|
|                      1 | TODO-PORT                |   done   |
|                      2 | TODO_DBFILE              |   done   | 
|                      3 | week , month             |   done   |  
|                      5 | /api/tasks?search=       |   done   | 
|                      8 | autorization, Dockerfile |   done   | 

 * Dockerfile 
```bash
# 1 member
docker build --tag scheduler:v1.0.0 .
docker run -d -p 8000:8000 scheduler:v1.0.0

# or 
# 2 member with compose.yaml
docker-compose --env-file ./init/.env up -d
 ```

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
├── internal
|   ├── database 
|   │   ├── database.go    // init for *sql.DB
|   │   ├── query.go       // SQL query for model
|   │   ├── schema.go      // SQL tables
|   │   └── transaction.go // *sql.DB, *sql.TX
|   ├── model              
|   │   ├──── login.go    
|   │   └──── task.go     
|   ├── server  
|   │   └──── server.go   // init for http.Server
|   ├── servises           
|   │   ├── autorization            // midlweare fucntion
|   │   │   └──── autorization.go   
|   │   ├── deserializer            // rules for get object from Request  
|   │   │   ├──── logindecode.go   
|   │   │   └──── taskdecode.go   
|   │   ├── serializer              // response computing & format
|   │   │   ├──── loginencode.go   
|   │   │   └──── taskencode.go   
|   │   ├── nextdate.go       // algorithnm for compute date in task
|   │   └── taskproperty.go   // rules for find task list   
|   └── transport   
|       ├── handler.go     // routes group
|       ├── route.go       
|       └── transport.go   // wrapper to '*http.ServeMux' and 'server.Srv'(server/server.go)
├── pkg/common
│   └──── common.go        // tools function
├── tests
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
main.go
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

#### For get .env file in main.go and test files
```bash
go get github.com/joho/godotenv # v1.5.1
```
#### Authorization use jwt.Token 
```bash
go get github.com/golang-jwt/jwt/v5 # v5.2.2
 ```
#### Tests
Some directories have file **coverage**, go to the folder you are interested in and call:
```bash
go tool cover -html=coverage
```
We can criate **coverage**:
```bash
# . - path
go test . -coverprofile=coverage.out
```
All test call
```bash
# for use ./test need run serve
go run main.go
# then 
go test ./...
```






p.s. Thanks for your time:)
















 



 



