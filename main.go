package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Section map[string]string
type AppInfo struct {
	Version string `json:"version"`
}

var appPathMapping = map[string]string{
	"app1": "appinfo.json",
	"app2": "app2/appinfo.json",
	// add more mappings here
}

func parseIniFile(path string) (map[string]Section, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]Section)
	section := "default"
	data[section] = make(Section)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			data[section] = make(Section)
			continue
		}

		if eqIndex := strings.Index(line, "="); eqIndex != -1 {
			key := strings.TrimSpace(line[:eqIndex])
			value := strings.TrimSpace(line[eqIndex+1:])
			data[section][key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func handleConfigs(w http.ResponseWriter, r *http.Request) {
	data, err := parseIniFile("test.ini")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.Query().Get("app_name")
	appPath, ok := appPathMapping[appName]
	if !ok {
		http.Error(w, "Invalid app name", http.StatusBadRequest)
		return
	}

	jsonFile, err := os.Open(appPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var appInfo AppInfo
	json.Unmarshal(byteValue, &appInfo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appInfo.Version)
}

func main() {
	http.HandleFunc("/configs", handleConfigs)
	http.HandleFunc("/version", handleVersion)
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
