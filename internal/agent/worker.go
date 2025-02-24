package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TimofeySar/ya_go_calculate.go/internal/calculation"
)

func Run(power int) {
	for i := 0; i < power; i++ {
		go worker()
	}
	select {}
}

func worker() {
	client := &http.Client{}
	for {
		resp, err := client.Get("http://localhost:8080/internal/task")
		if err != nil {
			fmt.Println("Error fetching task:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Println("Unexpected status:", resp.Status)
			time.Sleep(1 * time.Second)
			continue
		}

		var taskResp struct {
			Task calculation.Task `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
			resp.Body.Close()
			fmt.Println("Error decoding task:", err)
			continue
		}
		resp.Body.Close()

		task := taskResp.Task
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		var result float64
		switch task.Operation {
		case "+":
			result = task.Arg1 + task.Arg2
		case "-":
			result = task.Arg1 - task.Arg2
		case "*":
			result = task.Arg1 * task.Arg2
		case "/":
			if task.Arg2 == 0 {
				fmt.Println("Division by zero")
				continue
			}
			result = task.Arg1 / task.Arg2
		default:
			fmt.Println("Unknown operation:", task.Operation)
			continue
		}

		resultReq := struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}{
			ID:     task.ID,
			Result: result,
		}
		body, _ := json.Marshal(resultReq)
		resp, err = client.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error sending result:", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to send result, status:", resp.Status)
		}
	}
}
