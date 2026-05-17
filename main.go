package main

import (
	"basic-go/infra"
	"basic-go/src/controllers"
	"basic-go/src/repositories"
	"basic-go/src/services"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EchoRequest struct {
	Message string `json:"message"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title     string `json:"title"`
	Completed *bool  `json:"completed"`
}

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID = 1

func getTodoById(id int) *Todo {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i]
		}
	}
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req EchoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(req)
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	newTodo := Todo{ID: nextID, Title: req.Title}
	todos = append(todos, newTodo)
	nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newTodo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	todo := getTodoById(todoID)
	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if req.Completed == nil {
		http.Error(w, "Completed is required", http.StatusBadRequest)
		return
	}

	todo := getTodoById(todoID)
	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	updated := Todo{ID: todo.ID, Title: req.Title, Completed: *req.Completed}
	for i, t := range todos {
		if t.ID == todoID {
			todos[i] = updated
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	index := slices.IndexFunc(todos, func(t Todo) bool {
		return t.ID == todoID
	})
	if index == -1 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos = slices.Delete(todos, index, index+1)
	// デフォルトは200になるので200以外の場合は明示的にステータスを設定する必要あり
	w.WriteHeader(http.StatusNoContent)
}

// func main() {
// 	// Go 1.22から "メソッド /パス" の形式で登録できるようになり、マッチしないメソッドは自動で 405 Method Not Allowed を返す
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("GET /health", healthHandler)
// 	mux.HandleFunc("POST /echo", echoHandler)
// 	mux.HandleFunc("GET /todos", listTodosHandler)
// 	mux.HandleFunc("POST /todos", createTodoHandler)
// 	mux.HandleFunc("GET /todos/{id}", getTodoHandler)
// 	mux.HandleFunc("PUT /todos/{id}", updateTodoHandler)
// 	mux.HandleFunc("DELETE /todos/{id}", deleteTodoHandler)

// 	log.Println("server started at :8080")
// 	log.Fatal(http.ListenAndServe(":8080", mux))
// }

/***********************************************
ここからGin
************************************************/

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	// items := []models.Item{
	// 	{ID: 1, Name: "商品1", Price: 1000, Description: "説明1", SoldOut: false},
	// 	{ID: 2, Name: "商品2", Price: 2000, Description: "説明2", SoldOut: true},
	// 	{ID: 3, Name: "商品3", Price: 3000, Description: "説明3", SoldOut: false},
	// }

	// itemRepository := repositories.NewItemMemoryRepository(items)
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	r := gin.Default()
	r.GET("/items", itemController.FindAll)
	r.GET("/items/:id", itemController.FindById)
	r.POST("/items", itemController.Create)
	r.PUT("/items/:id", itemController.Update)
	r.DELETE("/items/:id", itemController.Delete)
	r.Run("localhost:8080") // デフォルトで0.0.0.0:8080でリッスンします
}
