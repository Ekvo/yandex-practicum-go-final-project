// This application: gets, saves, reads, updates and deletes a Task object, using the following patterns and tools.

// IMPORTANT: NOT TIME READ DOCS
//				 1. URL to start                            - http://127.0.0.1:8000 or http://127.0.0.1:8000/login.html
//           2. password to application is              - qwert12345
//           3. cookie life                             - 7 days     after autorization
//           4. test/settings.go
//           4.1. time life of Token in tests/settings.go - 7 days     start (~ 11.04.25 10:00)
//                after 7 days need set new. -> look bellow

// REST API property and features.

// Design patterns
/*
* CRUD
* Data Transfer Object (DTO)
* SOLID
 */

// DataBase    - SQLite
// multiplexer - ServerMUx

// package main ~> ./main.go
// logic of application

// package model ~> ../internal/model
/*
 - model.go
describe property of Task - object stored in the database
 * struct - TaskModel
 * 4 interface     - object maintenance
 * func UpdateDate - find next date for Task - use selected algorithm
 ---------------------------------------------------------------------------
describe property of Login
 * struct - LoginModel
 * func ValidPassword - compare password with TODO_PASSWORD (./init/.env)
*/

// packege server ~> ../internal/server
// rules for use http.Server in application
/*
 - server.go
 * struct - Srv                   - contain http.Server
 * func   - InitSRV               - get property from .env for initialize http.Serve
 * func   - ListenAndServeAndShut - property of connect and shut http.Server
*/

// packege servises ~> ../internal/servises
/*
logic of autorization
 - /autorization/autorization.go
 * func - AuthZ - midlweare function
 check autorization if exist TODO_PASSWORD (./init/.env)
 ---------------------------------------------------------------------------
 - /deserializer/logindecode.go
 * struct - LoginDecode   - create LoginModel from Request
 * func   - NewLoginDecode
 * func   - Model         - return LoginModel from LoginDecode
 * func   - Decode        - parse LoginDecode and create LoginModel
 ---------------------------------------------------------------------------
 - /deserializer/taskdecode.go
 * struct - TaskDecode    - create TaskModel from Request
 * func   - NewTaskDecode
 * func   - Model         - return TaskModel from LoginDecode
 * func   - Decode        - parse TaskDecode and create TaskModel
 * func   - executeDate   - rules for find 'data' when create new Task
 ---------------------------------------------------------------------------
 - /serializer/loginencode.go
 * struct - TokenResponse - return jwt.Token after SignedString -> (TokenGenerator) look (pkg/common/common.go)
 * struct - TokenEncode   - contain data for TokenResponse
 * func   - Response      - member of TokenEncode create TokenResponse
 ---------------------------------------------------------------------------
 - /serializer/taskencode.go
 * struct - TaskResponse   - create body from TaskModel for ResponseWriter
 * struct - TaskEncode     - contain data for TaskResponse
 * func   - Response       - member of TokenEncode create TaskResponse
 * struct - TaskListEncode - contain array of TaskModel for []TaskResponse
 * func   - Response       - member of TaskListEncode create array of TaskResponse
 ---------------------------------------------------------------------------
 - nextdate.go
algorithm for finding next date for TaskModel
 * func - NextDate        - main function for algorithm (find flag and call selected function by flag)
 * func - nextDateByDay   - create date by day(s) (UNIX - method)
 * func - numberOfDays    - find count of days for func nextDateByDay
 * func - nextDateByWeek  - date by day of week
 * func - daysOfWeek      - find UNIQUE days for func nextDateByWeek
 * func - nextDateByMonth - date by 1. number of day or 2.number of month(s) with number of day(s)
 * func - monthAndDay     - find days and month for func nextDateByMonth
 ---------------------------------------------------------------------------
 - taskproperty.go
rules for find 'model.TaskModel' array in database
 * struct - TaskProperty    - characteristics of Task(s) and lenght array []TaskModel from query (database)
 * func   - NewTaskProperty - call parseProperty setLimit
 * func   - setLimit        - member TaskProperty - find limit        (LIMIT)
 * func   - parseProperty   - member TaskProperty - find word or date (WHERE)
 * func   - IsDate          - member TaskProperty
 * func   - IsWord          - member TaskProperty
 * func   - PassDate        - member TaskProperty
 * func   - PassWord        - member TaskProperty
 * func   - PassLimite      - member TaskProperty
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
 ---------------------------------------------------------------------------
 - database.go
wrapper to 'dbTX' from transaction.go
 * struct - Source        - contain dbTX
 * func   - NewSource
 * func   - Init function - get property from .env for initialize database
 ---------------------------------------------------------------------------
 - schema.go
 * table(s) for database in format string
 ---------------------------------------------------------------------------
 - query.go
 * describe logic of interfaces Task (look: package model ~> ../internal/model/task.go)
*/

// packege transport ~> ../internal/transport
/*
 - query.go
 * describe logic of interfaces Task (look: package model ~> ../internal/model/task.go)
 ---------------------------------------------------------------------------
 - transport.go
wrapper to '*http.ServeMux' and 'server.Srv'
 * struct - Transport    - contain ptr of http.ServeMux and object server.Srv -> look (/internal/server/server.go)
 * func   - NewTransport
 * func   - Routes       - member Transport
 * func   - Run          - member Transport
 ---------------------------------------------------------------------------
 - route.go
describe application handlers
 ---------------------------------------------------------------------------
 - handler.go
rules for create route group
 * struct    - HandlerModel - empty struct
 * func      - NewHandlerModel
 * interface - multiTask   - contain all interfaces of TaskModel
 * func      - apiRoutes   - member HandlerModel - create group for address path /api/
*/

// packege common ~> ../pkg/common
/*
 - common.go
universal utilities
 * type      - Message -> (map)   - contain body for Response
 * func      - String             - member Message - formate message
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
jwt.Token - github.com/golang-jwt/jwt/v5
 * var       - SecretKey             - for use in 'jwtToken.SignedString'
 * func      - TokenGenerator        - generate token with (jwt.MapClaims) - content and time exploration
 * func      - ReceiveValueFromToken - chech time exploration and return value form token by key
*/

/*
 - Create new Token to test/settings.go
 ---------------------------------------------------------------------------
 1. run server (go run main.go)
 2. in other terminal:
curl -X POST -H "Content-Type: application/json" -d '{"password":"qwert12345"}' http://localhost:8000/api/signin

you should get a message like:
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250ZW50IjoiVGFzayBBY2Nlc3MiLCJleHBsb3JhdGlvbiI6MTc0NTEzODQ0M30.NGTz_-RhkPEEZZJ5uIku4DPtC0pCXQ6fjfrJCdM7M80"}

 3. copy jwtLINE from ("token":"jwtLINE") and set in test/settings.go -> Token = `jwtLINE`
*/

package docs
