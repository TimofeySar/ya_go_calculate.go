package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TimofeySar/ya_go_calculate.go/internal/orchestrator"
)

func main() {
	srv := orchestrator.NewServer()
	fmt.Println("Orchestrator running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
