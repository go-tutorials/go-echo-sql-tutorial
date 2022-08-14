# go-echo-postgresql-rest-api
an restful api service built with golang + echo framework + elephantSQL (postgresql cloud db)

## How to run
#### Clone the repository
```shell
git clone https://github.com/go-tutorials/go-echo-postgresql-rest-api.git
```

#### To run the application
```shell
go run main.go
```

## Database
#### How to create a cloud postgreSQL database with ElephantSQL
- create database document: https://www.elephantsql.com/docs/ (Note: use file data.sql (data -> data.sql) to get the sql query to create new database) 
- get URI link for the app to connect: choose the instance database you want to use after created in ElephantSQL, click "Details" option in the left column. Then copy the URL in the "Details" section
- use URI in code: configs -> config.yml. Then paste the URI above to the sql.data_source_name

## Echo Framework
#### get Echo Framework
```shell
go get github.com/labstack/echo/v4
go get github.com/labstack/echo/v4/middleware
```

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```

## API design for users
#### *Resource:* users

### Get all users
#### *Request:* GET /users
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "username": "peter.parker",
        "email": "peter.parker@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one user by id
#### *Request:* GET /users/:id
```shell
GET /users/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```

### Create a new user
#### *Request:* POST /users 
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: duplicate key, -1: error
```json
1
```

### Update one user by id
#### *Request:* PUT /users/:id
```shell
PUT /users/wolverine
```
```json
{
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

### Delete a new user by id
#### *Request:* DELETE /users/:id
```shell
DELETE /users/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include HealthHandler, HealthChecker, SqlHealthChecker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): log and log middleware

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
    db, err := sql.Open(conf.Driver, conf.DataSourceName)
    if err != nil {
        return nil, err
    }

    sqlChecker := s.NewSqlHealthChecker(db)
    healthHandler := health.NewHealthHandler(sqlChecker)
```

To handler routing
```go
    e := echo.New()	
    e.GET("/health", app.HealthHandler.Check)
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package main

import "github.com/core-go/config"

type Root struct {
    DB DatabaseConfig `mapstructure:"db"`
}

type DatabaseConfig struct {
    Driver         string `mapstructure:"driver"`
    DataSourceName string `mapstructure:"data_source_name"`
}

func main() {
    var conf Root
    err := config.Load(&conf, "configs/config")
    if err != nil {
        panic(err)
    }
}
```

### core-go/log *&* core-go/middleware
```go
import (
	"context"
	"fmt"

	"github.com/core-go/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-service/internal/app"
)

func main() {
	var conf app.Root

	e := echo.New()

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("Request Body: %v\n", string(reqBody))
		fmt.Printf("Response Body: %v\n", string(resBody))
		fmt.Printf("----------------------------------------\n")
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error}, message=${message}\n",
	}))
	e.Use(middleware.Recover())
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```

