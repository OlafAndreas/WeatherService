package main

import (
	"net/http"
	"fmt"
	"log"
	"time"
	"strings"
	"io/ioutil"
	xj "github.com/oisann/goxml2json"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"encoding/json"
)

func main() {

	setupCacheStorage()

	http.HandleFunc("/", api)
	logError(http.ListenAndServe(":8080", nil))
}

func api(w http.ResponseWriter, r *http.Request) {

	if path := r.URL.String(); path != "/favicon.ico" {

		// Split the string so we can retrieve the format in the last component
		components := strings.Split(path, ".")

		// Retrieve the format,
		// this will be used to transform the responses from YR.
		format := components[len(components)-1]

		w.Header().Set("Content-Type", "application/" + format)

		fmt.Fprintf(w, get(path, format))
	}
}

func print(text string) {
	log.Printf(text)
}

func get(relativeUrl string, format string) string {

	// All urls should be in lowercase,
	// this is to prevent redundancy in the cache DB.
	lowercased := strings.ToLower(relativeUrl)

	// Removing the user defined format in URL.
	// All requests to YR should be appended with xml.
	trimmed := strings.Trim(lowercased, format)

	url := "https://www.yr.no" + trimmed + "xml"

	// Checking the cache for a response
	cachedResponse := getCachedResponse(url)

	if len(cachedResponse) > 0 {
		print("Return cached response for: " + url)
		return formattedResponse(cachedResponse, format)
	}

	resp, error := http.Get(url)
	logError(error)

	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	logError(error)

	xml := string(body)

	// Here we should cache the data
	cacheResponse(url, xml)

	return formattedResponse(xml, format)
}

func formattedResponse(resp string, format string) string {

	if format == "json" {
		return xmlToJSON(resp)
	}

	return resp
}

func xmlToJSON(rawXml string) string {

	xml := strings.NewReader(rawXml)
  	json, error := xj.Convert(xml)
  	logError(error)

	return json.String()
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func database() *sql.DB {

	database, error := sql.Open("sqlite3", "./yrservicecache.db")
	logError(error)

	return database
}

func setupCacheStorage() {

	statement, _ := database().Prepare("CREATE TABLE IF NOT EXISTS cachedresponses (id INTEGER PRIMARY KEY, url TEXT, json TEXT, nextupdate DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL)")

	_, error := statement.Exec()
	logError(error)
}

func cacheResponse(url string, xml string) {
	print("Cache response for: " + url)

	var data map[string]interface{}
	err := json.Unmarshal([]byte(xmlToJSON(xml)), &data)
	logError(err)

	weatherdata, ok := data["weatherdata"].(map[string]interface{})

	if !ok {
		print("Failed to get weatherdata for caching.")
	}

	meta, ok := weatherdata["meta"].(map[string]interface{})
	if !ok {
		print("Failed to get meta data for caching.")
	}

	nextupdate, ok := meta["nextupdate"]
	if !ok {
		print("Failed to get nextUpdate from meta data for caching.")
	}

	statement, prepError := database().Prepare("INSERT INTO cachedresponses (url, json, nextupdate) VALUES (?, ?, ?)")
	logError(prepError)

	_, execError := statement.Exec(url, xml, nextupdate)
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

// maxAge is measured in minutes
func getCachedResponse(url string) string {
	print("Check cache for: " + url)
	rows, queryError := database().Query("SELECT json, nextupdate FROM cachedresponses WHERE url=?", url)
	logError(queryError)

	defer rows.Close()

	var json string
	var nextupdate time.Time
	for rows.Next() {
		scanError := rows.Scan(&json, &nextupdate)
		logError(scanError)
	}

	logError(rows.Err())

	if time.Since(nextupdate).Minutes() >= 0 {
		removeOldCache(url)
		return ""
	}

	return json
}
