package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(todos) == 0 {
		message := map[string]string{"message": "No todos available"}
		err := json.NewEncoder(w).Encode(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func getTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, v := range todos {
		if v.ID == params["id"] {
			json.NewEncoder(w).Encode(v)
			return
		}
	}
	message := map[string]string{"message": "Todo not found"}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(message)
}

func deleteTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	params := mux.Vars(r)
	for i, v := range todos {
		if v.ID == params["id"] {
			todos = append(todos[:i], todos[i+1:]...)
			message := map[string]string{"message": "todo has been deleted successfully"}
			json.NewEncoder(w).Encode(message)
			break
		}
		json.NewEncoder(w).Encode(todos)
	}
}

var idCounter = 3

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todo.ID = strconv.Itoa(idCounter)
	idCounter++
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

func editTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, v := range todos {
		if v.ID == params["id"] {
			todos = append(todos[:i], todos[i+1:]...)

			var todo Todo
			err := json.NewDecoder(r.Body).Decode(&todo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			todo.ID = params["id"]
			todos = append(todos, todo)
			json.NewEncoder(w).Encode(todo)
		}
	}
}

func main() {
	r := mux.NewRouter()

	todos = append(todos, Todo{"1", "eat", false})
	todos = append(todos, Todo{"2", "sleep", true})

	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todo/{id}", getTodoById).Methods("GET")
	r.HandleFunc("/todo/create", CreateTodo).Methods("POST")
	r.HandleFunc("/todo/edit/{id}", editTodoById).Methods("PUT")
	r.HandleFunc("/todo/delete/{id}", deleteTodoById).Methods("DELETE")
	fmt.Printf("Starting server at port 8080.....\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
