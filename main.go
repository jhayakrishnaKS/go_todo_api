package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Todo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB

func initDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/todo_api_go?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	db.AutoMigrate(&Todo{})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo
	db.First(&todo, params["id"])
	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	db.Create(&todo)
	json.NewEncoder(w).Encode(todo)
	json.NewEncoder(w).Encode("Todo has been created successfully")

}

func UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo
	db.First(&todo, params["id"])
	json.NewDecoder(r.Body).Decode(&todo)
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
	json.NewEncoder(w).Encode("Todo has been updated successfully")

}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo
	db.First(&todo, params["id"])
	db.Delete(&todo)
	json.NewEncoder(w).Encode("Todo deleted successfully")
}

func main() {
	initDB()
	defer db.Statement.ReflectValue.Close()
	r := mux.NewRouter()
	r.HandleFunc("/todos", GetTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", GetTodoById).Methods("GET")
	r.HandleFunc("/todos", CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", UpdateTodoById).Methods("PUT")
	r.HandleFunc("/todos/{id}", DeleteTodoById).Methods("DELETE")

	fmt.Println("Starting server at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
