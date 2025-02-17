package main

import (
	"encoding/json"
	//"errors"
	//"fmt"
	"log"
	//"math/rand"
	"bytes"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result *float64 `json:"result,omitempty"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

var (
	expressions = make(map[string]*Expression)
	tasks       = make([]*Task, 0)
	mutex       sync.Mutex
)


func calculateHandler(w http.ResponseWriter, r *http.Request) {

}

func getExpressionsHandler(w http.ResponseWriter, r *http.Request) {

}

func getExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func submitResultHandler(w http.ResponseWriter, r *http.Request) {

}

func worker() {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Println("Failed to fetch task:", err)
			time.Sleep(time.Second)
			continue
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			log.Println("Invalid task response:", err)
			continue
		}

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		var result float64
		switch task.Operation {
		case "+": result = task.Arg1 + task.Arg2
		case "-": result = task.Arg1 - task.Arg2
		case "*": result = task.Arg1 * task.Arg2
		case "/": if task.Arg2 != 0 { result = task.Arg1 / task.Arg2 } else { log.Println("Division by zero") }
		}

		data, _ := json.Marshal(map[string]interface{}{"id": task.ID, "result": result})
		_, err = http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewReader(data))
		if err != nil {
			log.Println("Failed to submit result:", err)
		}
	}
}

func main() {
	workerCount, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if workerCount == 0 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		go worker()
	}

	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/api/v1/expressions", getExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", getExpressionByIDHandler)
	http.HandleFunc("/internal/task", getTaskHandler)
	http.HandleFunc("/internal/task", submitResultHandler)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
