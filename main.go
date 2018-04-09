package main

import (
	"net/http"
	"log"
	"time"
	"strings"
	"io/ioutil"
	xj "github.com/oisann/goxml2json"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

func main() {

	setupCacheStorage()

	get("/Trøndelag/Stjørdal/Stjørdal") // Fetched from server
	get("/Trøndelag/Stjørdal/Stjørdal") // Fetched from cache

	get("/Trøndelag/Verdal/Verdal")	// Fetched from server
	get("/Trøndelag/Verdal/Verdal") // Fetched from cache
}

func print(text string) {
	log.Printf(text)
}

func get(relativeUrl string) string {

	url := "http://www.yr.no/sted/Norge"+relativeUrl+"/varsel.xml"

	// Checking the cache for a response added in the last 10 minutes.
	cachedResponse := getCachedResponse(url, 10)

	if len(cachedResponse) > 0 {
		print("Return cached response for: " + url)
		return cachedResponse
	}

	resp, error := http.Get(url)
	logError(error)

	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	logError(error)

	json := xmlToJSON(string(body))

	// Here we should cache the data
	cacheResponse(url, json)

	return json
}

func xmlToJSON(rawXml string) string {

	xml := strings.NewReader(rawXml)
  	json, error := xj.Convert(xml)
  	logError(error)

	return json.String()
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func database() *sql.DB {

	database, error := sql.Open("sqlite3", "./yrservicecache.db")
	logError(error)

	return database
}

func setupCacheStorage() {

	statement, _ := database().Prepare("CREATE TABLE IF NOT EXISTS cachedresponses (id INTEGER PRIMARY KEY, url TEXT, json TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")

	_, error := statement.Exec()
	logError(error)
}

func cacheResponse(url string, json string) {
	print("Cache response for: " + url)
	statement, prepError := database().Prepare("INSERT INTO cachedresponses (url, json) VALUES (?, ?)")
	logError(prepError)

	_, execError := statement.Exec(url, json)
	logError(execError)
}

func removeOldCache(url string) {

	statement, prepError := database().Prepare("DELETE FROM cachedresponses WHERE url=?")
	logError(prepError)

	result, execError := statement.Exec(url)
	logError(execError)
	

	_, affectedError := result.RowsAffected()
	logError(affectedError)
}

// maxAge is measured in seconds
func getCachedResponse(url string, maxAge float64) string {
	print("Check for cached response for: " + url)
	rows, queryError := database().Query("SELECT json, timestamp FROM cachedresponses WHERE url=?", url)
	logError(queryError)

	defer rows.Close()

	var json string
	var timestamp time.Time
	for rows.Next() {
		scanError := rows.Scan(&json, &timestamp)
		logError(scanError)
	}

	logError(rows.Err())

	if time.Since(timestamp).Minutes() > maxAge {
		removeOldCache(url)
		return ""
	}

	return json
}