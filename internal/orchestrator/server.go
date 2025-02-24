package orchestrator

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/TimofeySar/ya_go_calculate.go/internal/calculation"
	"github.com/gorilla/mux"
)

type Server struct {
	expressions map[string]*Expression
	tasks       chan calculation.Task
	mu          sync.Mutex
}

func NewServer() *mux.Router {
	srv := &Server{
		expressions: make(map[string]*Expression),
		tasks:       make(chan calculation.Task, 100),
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", srv.handleCalculate).Methods("POST")
	r.HandleFunc("/api/v1/expressions", srv.handleListExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", srv.handleGetExpression).Methods("GET")
	r.HandleFunc("/internal/task", srv.handleGetTask).Methods("GET")
	r.HandleFunc("/internal/task", srv.handleReceiveResult).Methods("POST")
	return r
}

func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}
	if req.Expression == "" || strings.ContainsAny(req.Expression, "!@#$%^&") {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	id := generateID()
	expr := NewExpression(id, req.Expression)
	s.mu.Lock()
	s.expressions[id] = expr
	s.mu.Unlock()

	go expr.Start(s.tasks)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleListExpressions(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp := struct {
		Expressions []Expression `json:"expressions"`
	}{}
	for _, expr := range s.expressions {
		resp.Expressions = append(resp.Expressions, *expr)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.mu.Lock()
	expr, exists := s.expressions[id]
	s.mu.Unlock()
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]Expression{"expression": *expr}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-s.tasks:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]calculation.Task{"task": task}); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func (s *Server) handleReceiveResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}
	// Добавляем валидацию id
	if req.ID == "" {
		http.Error(w, "Invalid task ID", http.StatusUnprocessableEntity)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, expr := range s.expressions {
		if expr.UpdateTaskResult(req.ID, req.Result) {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
