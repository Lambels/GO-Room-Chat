[![Version](https://img.shields.io/badge/goversion-1.16.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>

# Go Room Chat

Room chat implemented with go using Gorilla mux and Gorilla websockets

## Build

Mac/Linux:

~~~
go build -o GoRoomChat ./cmd/web/*.go
~~~

Windows:

~~~
go build -o GoRoomChat.exe ./cmd/web/.
~~~

For particular platform:

~~~
env GOOS=linux GOARCH=amd64 go build -o GoRoomChat ./cmd/web/*.go
~~~

## Requirements

GoRoomChat Requires:
- A running mysql database

## Run

Make sure the database is running.

(If you are running on localhost with a database with no password default config should be ok)
Run with flags:

~~~
./GoRoomChat \
-dbUser='root' \
-dbName='testdb' \
-dbPass='' \
-dbPort='3306' \
-dbHost='localhost' \
-stName='sessions' \
-port=':8080'
~~~

## All Flags

~~~
./GoRoomChat -help
Usage of GoRoomChat:
-dbHost string
    Database host (default "localhost")
-dbName string
    Database name (default "testdb")
-dbPass string
    Database password
-dbPort string
    Database port (default "3306")
-dbUser string
    Databse user (default "root")
-port string
    Port on wich the server listens on (default ":8080")
-stName string
    Sessions table name (default "sessions")
~~~