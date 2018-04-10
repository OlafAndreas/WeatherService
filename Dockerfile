FROM golang:1.10
WORKDIR /go/src/github.com/olafandreas/yrservice
RUN go get -u github.com/mattn/go-sqlite3
RUN go get -u github.com/oisann/goxml2json
COPY . .
RUN go build -o YRService .
CMD ./YRService