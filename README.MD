# WeatherService
#### A simple Go application that communicates with the YR API.

Example API Url: [http://www.yr.no/sted/Norge/Trøndelag/Verdal/Verdal/varsel.xml](http://www.yr.no/sted/Norge/Tr%C3%B8ndelag/Verdal/Verdal/varsel.xml)

## Installation

For instructions on installation of Go, please refer to [https://golang.org/dl/](https://golang.org/dl/).

##### Install dependencies

###### Using dep
This project uses dep, the prototype dependency manager for Go.
Please refer to [https://github.com/golang/dep](https://github.com/golang/dep) for installation.

Run the following command inside the WeatherService directory to install the dependencies.
```
dep ensure
```
And you're good to go! :)

#### or

###### Using `go get`
Packages can be insatlled manually by running the following commands, I recommend using dep since these commands will install the latest packages available, which might not work in the future.
```
go get github.com/mattn/go-sqlite3
go get github.com/oisann/goxml2json
```
And you might be good to go! :)

## Build & Run
WeatherService will start a web service on port 8080.
```
go build

./WeatherService
```

## Usage
Retrieve data as XML:
```
http://127.0.0.1:8080/sted/norge/trøndelag/verdal/vinne/varsel.xml
```
Retrieve data as JSON:
```
http://127.0.0.1:8080/sted/norge/trøndelag/verdal/vinne/varsel.json
```

Both responses would be cached using `http://127.0.0.1:8080/sted/norge/oslo/oslo/kjelsås/varsel.` as the key, only the first response will be fetched from the yr api and cached, any requests on the same url will return the cached response in the specific format. The cache database is purged every 5 minutes to remove any expired caches.

#### Practical information
This application was built using `go1.10.1 darwin/amd64`
