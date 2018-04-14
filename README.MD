# WeatherService
#### A simple Go application that communicates with the YR API.

Example API Url: [http://www.yr.no/sted/Norge/Trøndelag/Verdal/Verdal/varsel.xml](http://www.yr.no/sted/Norge/Tr%C3%B8ndelag/Verdal/Verdal/varsel.xml)

## Installation

For instructions on installation of Go, please refer to [https://golang.org/dl/](https://golang.org/dl/).

Installing dependencies using [dep](https://github.com/golang/dep).
The following command will take care of installing any dependencies required
by the WeatherService app.
```
dep ensure
```

## Build & Run
```
go build

./WeatherService // Starts a web service on port 8080 by default
```

## Usage
Retrieve data as XML:
```
http://127.0.0.1:8080/sted/norge/trøndelag/verdal/vinne/varsel.xml
```

Retrieve data as JSON:
```
http://127.0.0.1:8080/sted/norge/trøndelag/verdal/vinne/varsel.JSON
```

Both responses would be cached using `http://127.0.0.1:8080/sted/norge/oslo/oslo/kjelsås/varsel.` as the key, only the first response will be fetched from the yr api and cached, any requests on the same url will return the cached response in the specific format. The cache database is purged every 5 minutes to remove any expired caches.

#### Practical information
This application was built using `go1.10.1 darwin/amd64`
