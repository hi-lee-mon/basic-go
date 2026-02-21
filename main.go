package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type EchoRequest struct {
	Message string `json:"message"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// データストア（メモリ上）
var todos []Todo
var nextID int = 1

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

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
	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
