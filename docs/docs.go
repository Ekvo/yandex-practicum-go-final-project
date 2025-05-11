// This application: gets, saves, reads, updates and deletes a Task object, using the following patterns and tools.

// IMPORTANT: NOT TIME READ DOCS
// 1. URL to start                              - http://127.0.0.1:8000 or http://127.0.0.1:8000/login.html
// 2. password to application is                - qwert12345
// 3. cookie life                               - 7 days     after autorization
// 4. test/settings.go
// 4.1. time life of Token in tests/settings.go - 7 days     start (~ 20.04.25 10:00)
// after 7 days need set new. -> look bellow

//======================================================================================================
// REST API property and features.
//======================================================================================================

//======================================================================================================
// Design patterns
//======================================================================================================
/*
* CRUD
* Data Transfer Object (DTO)
* SOLID
 */

// DataBase    - SQLite
// multiplexer - ServerMUx
//======================================================================================================

// package main ~> .cmd/app/main.go
// logic of application

// package model ~> ../internal/app
/*
 - app.go
main object of application
 * struct    - Sheduler    - contain all interfaces of application
 * func      - NewSheduler
 ------------------------------------------------------------------------------------------------------
 - run.go
 * func - Run  - start and close of application
*/

// package config ~> ../internal/config
/*
 - config.go
rules for create config
 * struct - Config           - container for fields responsible for creating application service objects
 * func   - NewConfig
 * func   - loadConfig       - member of Config - read file and fill the config
 * func   - setConfig        - member of Config - set extension, ser Reader to 'viper.ReadConfig'
 * func   - ValidConfig      - member of Config - check of config on valid, call all 'validSomeFiled'
 * func   - validDataBase    - member of Config
 * func   - validServe       - member of Config
 * func   - validTask        - member of Config
 * func   - validPassword    - member of Config
 * func   - validJWT         - member of Config
 * func   - validPathOfFiles - member of Config
 ------------------------------------------------------------------------------------------------------
 - options.go
property of file for config onject
 * struct - options   - initial rules for creating a configuration
 * file   - parsePath - fiil fileds options (!without pathOfFile!) with parse 'options.pathOfFile'
*/

// package model ~> ../internal/model
/*
 - task.go
describes property of Task - object stored in the database
 * struct      - TaskModel
 * 4 interface - TaskModel object maintenance in repository
 ------------------------------------------------------------------------------------------------------
describe property of Login
 - login.go
 * struct    - LoginModel
 * interface - LoginRead  - check 'login' on exist and if exist check password
*/

// packege source ~> ../internal/database
// SQLite - modernc.org/sqlite"
/*
 - transaction.go
wrapper to '*sql.DB' and '*sql.Tx'
 * struct - dbTX        - contain ptrs of sql.DB and sql.Tx
 * func   - Transaction - member dbTX rules for create,update, datele source in databse
 Transaction - 1. create ptr of sql.TX (call (db *DB) BeginTx)
               2. set dbTX.TX (dbtx.Tx = tx)
               3. calls the function passed to it - 'execute'
					4. 'execute' should insert,update,delete only use 'dbtx.Tx.(type of query)'
 ------------------------------------------------------------------------------------------------------
 - database.go
wrapper to 'dbTX' from transaction.go
 * struct - Source        - contain dbTX
 * func   - NewSource
 * func   - Init function - get property from config for initialize database
 ------------------------------------------------------------------------------------------------------
 - schema.go
 * table(s) for database in format string
 ------------------------------------------------------------------------------------------------------
 - query.go
 * describe logic of interfaces Task (look: package model ~> ../internal/model/task.go)
*/

// package datauser ~> ../internal/datauser
// cashe store for user data
/*
 - datauser.go
 * struct - UserData      - contain password of user (hash)
 * func   - NewUserData   - create UserData, set password from config

implementatain of 'LoginRead interface' use object 'UserData'
 * func   - ValidLogin    - member of UserData - compare lines
 * func   - PasswordExist - member of UserData - compare lenght of password with 0
*/

// package jwtsign ~> ../internal/lib/jwtsign
// contain 'secretkey' for create, parse 'jwt.Token'
/*
 - jwtsign.go
 * var  - secretKey                - non-exported global variable
 * func - NewSecretKey             - set secretKey used config (used only in start application and if secretKey is empty)
 * func - TokenRetrieve            - wrapper for 'jwt.Parse'
 * func - TokenGenerator           - creates jwt.Token with 'exploration' and set string  by key 'content'
 * func - ReceiveValueFromToken[T] - generic - get value[T] from token by key 'jwt.MapClaims'
*/

// package nextdate ~> ../internal/lib/nextdate
// algorithm for finding next date for TaskModel
/*
 - nextdate.go
 * type - NextDateFunc    - func(now time.Time, dstart string, repeat string) (string, error) - describes logic of algorithm for find date in model.TaskModel
 * func - NextDate        - main function for algorithm (find flag and call selected function by flag)
 * func - nextDateByDay   - create date by day(s) (UNIX - method)
 * func - numberOfDays    - find count of days for func nextDateByDay
 * func - nextDateByWeek  - date by day of week
 * func - daysOfWeek      - find UNIQUE days for func nextDateByWeek
 * func - nextDateByMonth - date by 1. number of day or 2. number of month(s) with number of day(s)
 * func - monthAndDay     - find days and month for func nextDateByMonth
*/

// packege server ~> ../internal/server
// rules for use http.Server in application
/*
 - server.go
 * struct - Srv                   - contain http.Server
 * func   - InitSRV               - get property from config for initialize http.Serve
 * func   - ListenAndServeAndShut - property of connect and shut http.Server
*/

// packege servises ~> ../internal/servises
/*
 - services.go
describes all biz logic of application
create content for ResponseWriter
 * interface - TaskCreateCase
 \_ 'CreateTask' - take 'model.TaskModel' and return '*serializer.TaskIDResponse',error
 * interface - TaskReadCase
 |_ 'ReadTask'   - take 'uint' ID of task and return '*serializer.TaskResponse'
 \_ ReadTaskList - take '*entity.TaskProperty' see (../internal/services/entity/taskproperty.go) and return '*serializer.TaslListResponse',error
 * interface - TaskUpdateCase
 \_ 'UpdateTask' - take 'model.TaskModel' for update in some store and return only error
 * interface - TaskDeleteCase
 \_ 'DeleteTask' - take 'uint' for delete task by ID and return only error
 * interface - TaskDoneCase
 \_ 'DoneTask' - take 'uint' for update status(update or delete) task by ID and return only error
 * interface - LoginValidPasswordCase
 |_ 'CreateToken' - take 'model.LoginModel' for create 'jwt.Token' after return '*serializer.TokenResponse', error
 \_ 'UserExist'   - check login exist in application, return bool,error
 * interface - AutorizationCase
 \_'AuthZ' - take '*http.Request' and find user data
*/

// packege usecase ~> ../internal/servises/usecase
// implement of 'services.go' interfaces
/*
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 - authcase.go
 * interface - AuthService    - contain all biz logic interfaces Autorization
 * struct    - authService    - empty struct for implemet and use member of interface - AutorizationCase from (services.go)
 * func      - NewAuthService
 * func      - AuthZ          - describes biz logic of autorization
 checks 'token' in cookie by key 'token', after parse token and check for validity, find field in token by key 'content' after prints received line
 ------------------------------------------------------------------------------------------------------
 - logincase.go
 * interface - LoginService    - contain all business logic interfaces LoginCase
 * interface - MultiLogin      - all interfaces of 'model.LoginRead' work with store
 * struct    - loginService    - have a loginRepository logic work with MultiLogin
 * func      - NewLoginService
 * func      - UserExist       - describes biz logic of check on exist 'login(datauser/datauser.go)' in application
 * func      - CreateToken     - logic of create object(serializer.TokenResponse) with jwt.Token inside
 ------------------------------------------------------------------------------------------------------
 - taskcase.go
 * interface - TaskService - contain all business logic interfaces of all Task Case
 * interface - MultiTask   - all interfaces of 'model.TaskModel work with store
 * struct    - taskService
1. taskRepository logic -> work with MultiTask (internal/database/)
2. have a algorithm 'nextdate.NextDateFunc' (lib/nextdate/nextdate.go) for find next date of Task
 * func      - NewTaskService
 * func      - setNextDate         - get name of algorithm from config and return function of type 'nextdate.NextDateFunc'
 * func      - CreateTask          - logic of create task (more information insade package)
 * func      - executeDate         - finds date when a task was created or updated (details in package)
 * func      - ReadTask            - logic of read Task by ID from database and create object for Response
 * func      - UpdateTask          - rules for update Task By ID
 * func      - DeleteTask          - describes process of deletin task by ID from store
 * func      - DoneTask            - done task by ID
1. use func 'updateDateAfterDone' see bellow (more details in package)
2. if rules for repeat Task is empty - delete Task from store
3. othercase update task in database
 * func      - updateDateAfterDone - finds date when a task was done
 * func      - ReadTaskList        - create Task List for response, by rules:(*entity.TaskProperty) see (/service/entity/taskproperty.go)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/

//packege deserializer ~> ../internal/servises/deserializer
/*
 - taskdecode.go
 * struct - TaskDecode    - create TaskModel from Request
 * func   - NewTaskDecode
 * func   - Model         - return TaskModel from LoginDecode
 * func   - Decode        - parse TaskDecode and create TaskModel
 * func   - executeDate   - rules for find 'data' when create new Task
 ------------------------------------------------------------------------------------------------------
 - /deserializer/logindecode.go
 * struct - LoginDecode   - create LoginModel from Request
 * func   - NewLoginDecode
 * func   - Model         - return LoginModel from LoginDecode
 * func   - Decode        - parse LoginDecode and create LoginModel
*/

// package serializer ~>  ../internal/servises/serializer
// rules for create objects for Response
/*
 - loginencode.go
 * struct - TokenResponse - return jwtsign.Token after SignedString -> (TokenGenerator) look (pkg/common/common.go)
 * struct - TokenEncode   - contain data for TokenResponse
 * func   - Response      - member of TokenEncode create TokenResponse
 ------------------------------------------------------------------------------------------------------
 - taskencode.go
 * struct - TaskResponse     - object contain one Task for Response
 * struct - TaskEncode       - contain start data for TaskResponse
 * func   - Response         - member of TokenEncode create TaskResponse
 * struct - TaslListResponse - object contain array of Task for Response
 * struct - TaskListEncode   - contain array of TaskModel
 * func   - Response         - member of TaskListEncode create array of TaskResponse
 * struct - TaskIDResponse   - Task ID Transfer Rules
 * strcut - TaskIDEncode     - have a positive number of Task
 * func   - Response         - member of TaskIDEncode create TaskIDResponse
*/

// package entity ~> ../internal/services/entity
/*
 - taskproperty.go
rules for find 'model.TaskModel' array in database
 * struct - TaskProperty    - characteristics of Task(s) and lenght array []TaskModel from query (database)
 * func   - NewTaskProperty - call parseProperty setLimit
 * func   - setLimit        - member TaskProperty - find limit        (LIMIT)
 * func   - parseProperty   - member TaskProperty - find word or date (WHERE)
 * func   - IsDate          - member TaskProperty - find by date
 * func   - IsWord          - member TaskProperty - find by word
 * func   - PassDate        - member TaskProperty
 * func   - PassWord        - member TaskProperty
 * func   - PassLimite      - member TaskProperty
*/

// packege transport ~> ../internal/transport
/*
 - transport.go
wrapper to '*http.ServeMux' and 'server.Srv'
 * struct    - Transport    - contain ptr of http.ServeMux and object server.Srv -> look (/internal/server/server.go)
 * func      - NewTransport
 * interface - shedulerCase   - contain all interfaces from services
pass the Route function as an argument, set a specific interface in each handler in route.go
 * func   - Routes       - member Transport
 * func   - Run          - member Transport
 ------------------------------------------------------------------------------------------------------
 - middlweare.go
describes middlweare functions
 * interface - rulesForAuthZ - rules for check autorization
 * func      - AuthZ         - take next('http.HandlerFunc') and check with help 'rulesForAuthZ'
1. if login !exist -> call next
2. othercase check password -> cal next (details inside file - middlweare.go)
 ------------------------------------------------------------------------------------------------------
 - route.go
describe application handlers
 ------------------------------------------------------------------------------------------------------
 - handler.go
rules for create route group
 * struct    - HandlerModel - empty struct
 * func      - NewHandlerModel
 * func      - apiRoutes      - member HandlerModel - create group for address path /api/
*/

// packege common ~> ../pkg/common
/*
 - common.go
universal utilities
 * type      - Message -> (map)   - contain body for Response
 * func      - String             - member Message - formate message
 * struct    - MessageError       - format for error message
 * func      - NewError
 * interface - ScanSQL            - contain Scan(dest ...any) error
 * func      - CreatePathWithFile - take the path with the file name and add it to the current path,
                                    then create folders with the file if needed
 * func      - Abs                - absolute value: |-(7+9)| = 16
 * func      - DecodeJSON         - common rules for every objects from Request
 * func      - EncodeJSON         - rules for create body to Response
 * func      - BeginningOfMonth   - work with time.Time
 * func      - ReduceTimeToDay    - time.Time
 * func      - HashData           - hashes some string string
 * func      - ReadCookie         - find cookie.Value by key
 * func      - CleanCookie        - end all cookies life
*/

/*
 - Create new Token to test/settings.go
 ------------------------------------------------------------------------------------------------------
 1. run server (go run cmd/app/main.go)
 2. in other terminal:
curl -X POST -H "Content-Type: application/json" -d '{"password":"qwert12345"}' http://localhost:8000/api/signin

you should get a message like:
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250ZW50IjoiVGFzayBBY2Nlc3MiLCJleHBsb3JhdGlvbiI6MTc0NTEzODQ0M30.NGTz_-RhkPEEZZJ5uIku4DPtC0pCXQ6fjfrJCdM7M80"}

 3. copy jwtLINE from ("token":"jwtLINE") and set in test/settings.go -> Token = `jwtLINE`
*/

package docs
