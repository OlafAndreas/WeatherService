# YRService
#### A simple Go application that communicates with the YR API.

Example API Url: [http://www.yr.no/sted/Norge/Trøndelag/Verdal/Verdal/varsel.xml](http://www.yr.no/sted/Norge/Tr%C3%B8ndelag/Verdal/Verdal/varsel.xml)

## Installation
Install these two packages
```
go get github.com/mattn/go-sqlite3
go get -u github.com/oisann/goxml2json
```

## Build & Run
```
go build

./YRService
```

This application was built using `go1.10.1 darwin/amd64`