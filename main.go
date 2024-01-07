package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// ApmRequest represents the expected JSON format for the /APM_switch endpoint.
type ApmRequest struct {
	Switch string `json:"switch"`
	Fab    string `json:"fab"`
}

func apmSwitchHandler(w http.ResponseWriter, r *http.Request) {
	var request ApmRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Assuming apm.ps1 is in the same directory as the executable.
	scriptPath, err := filepath.Abs("D:/Documents/VScode/VM_Service/VM_Service/apm.ps1")
	if err != nil {
		log.Println("Error finding script path:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Printf(scriptPath)

	// Constructing the PowerShell command.
	command := exec.Command("powershell", "-File", scriptPath, request.Switch, request.Fab)
	output, err := command.CombinedOutput()
	if err != nil {
		log.Println("Error running PowerShell script:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := strings.TrimSpace(string(output))
	fmt.Fprintf(w, "Script Output:\n%s", response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/APM_switch", apmSwitchHandler).Methods("POST")

	port := 8080 // You can change the port as needed
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
