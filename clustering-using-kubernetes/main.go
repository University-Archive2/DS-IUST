package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func instanceIDHandler(w http.ResponseWriter, r *http.Request) {
	podName := os.Getenv("POD_NAME")

	if podName == "" {
		http.Error(w, "Pod Name not found", http.StatusInternalServerError)
		return
	}

	parts := strings.Split(podName, "-")
	podID := parts[len(parts)-1]

	fmt.Fprintf(w, "Pod id: %s", podID)
}

func main() {
	http.HandleFunc("/instance-id", instanceIDHandler)
	fmt.Println("Server is starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
