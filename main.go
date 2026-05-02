package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// データストア（メモリ上）
var todos []Todo
var nextID int = 1

func getTodoById(id int) *Todo {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i]
		}
	}
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	http.HandleFunc("/health", healthHandler)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		var req EchoRequest
		// bodyの内容をreqの構造体に代入する。値を変更するときはGoではポインタを使う必要がある。
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(req)
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		allowed := map[string]struct{}{
			http.MethodGet:  {},
			http.MethodPost: {},
		}
		if _, ok := allowed[r.Method]; !ok {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			// リクエストの内容を構造体にデコード
			var req CreateTodoRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			// バリデーション
			if req.Title == "" {
				http.Error(w, "Title is required", http.StatusBadRequest)
				return
			}

			// レスポンス作成
			newTodo := Todo{
				ID:        nextID,
				Title:     req.Title,
				Completed: false,
			}

			// ストア更新
			todos = append(todos, newTodo)
			nextID++

			// ステータスコード201を返す
			w.WriteHeader(http.StatusCreated)
			// レスポンス返却
			_ = json.NewEncoder(w).Encode(newTodo)
		}

		if r.Method == http.MethodGet {
			_ = json.NewEncoder(w).Encode(todos)
		}

	})

	http.HandleFunc("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		allowed := map[string]struct{}{
			http.MethodGet: {},
			http.MethodPut: {},
		}
		if _, ok := allowed[r.Method]; !ok {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// パスパラの取得と検証
		idStr := r.PathValue("id")
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		todoID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID must be a number", http.StatusBadRequest)
			return
		}
		/*
			Get
		*/
		if r.Method == http.MethodGet {
			todo := getTodoById(todoID)

			if todo == nil {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(todo)
		}

		/*
			Put
		*/
		if r.Method == http.MethodPut {
			// リクエストの内容を構造体にデコード
			var req UpdateTodoRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			// バリデーション
			if req.Title == "" {
				http.Error(w, "Title is required", http.StatusBadRequest)
				return
			}

			if req.Completed == nil {
				http.Error(w, "Completed is required", http.StatusBadRequest)
				return
			}

			// 検索
			todo := getTodoById(todoID)

			if todo == nil {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			}

			updatedTodo := Todo{
				ID:        todo.ID,
				Title:     req.Title,
				Completed: *req.Completed,
			}

			// ストア更新
			for i, t := range todos {
				if t.ID == todoID {
					todos[i] = updatedTodo
					break
				}
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(updatedTodo)
		}
	})
	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
