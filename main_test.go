package main

import (
	"basic-go/infra"
	"basic-go/src/dto"
	"basic-go/src/models"
	"basic-go/src/services"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestMainとすることでvitestのbeforeAllの動きをする
func TestMain(m *testing.M) {
	// テストが実行される前に一度だけ必ず.env.testファイルを読み込む
	if err := godotenv.Load("./.env.test"); err != nil {
		log.Fatal("Error loading .env.test file")
	}

	// 全テストを実行
	code := m.Run()
	// テストが完了したら終了
	os.Exit(code)
}

func setupTestData(db *gorm.DB) {
	items := []models.Item{
		{Name: "Test Item 1", Description: "This is a test item 1", Price: 10.0, SoldOut: false, UserId: 1},
		{Name: "Test Item 2", Description: "This is a test item 2", Price: 20.0, SoldOut: true, UserId: 2},
	}

	users := []models.User{
		{Email: "test1@example.com", Password: "password0001"},
		{Email: "test2@example.com", Password: "password0002"},
	}

	for _, item := range items {
		db.Create(&item)
	}

	for _, user := range users {
		db.Create(&user)
	}
}

/**
* テスト用のテーブル作成、データ投入、ルーターのセットアップを行う関数
* これを各テスト関数の冒頭で呼び出すことで、テストごとにクリーンな状態でテストを実行できるようになる
**/
func setup() *gin.Engine {
	// DBの初期化(これが実行される前にTestMainで.env.testが読み込まれているため、テスト用のDBに接続される)
	mockDb := infra.SetupDB()
	// マイグレーションを実行してテスト用のテーブルを作成
	mockDb.AutoMigrate(&models.Item{}, &models.User{})

	// 作成したテーブルにデータをインサート
	setupTestData(mockDb)

	// テスト用のルーターをセットアップ
	router := setupRouter(mockDb)

	return router
}

func TestFindAll(t *testing.T) {
	/**
	* Arrange
	**/
	// テスト用のルーターをセットアップ
	router := setup()
	// リクエストとレスポンスの記録
	w := httptest.NewRecorder()
	// リクエスト作成
	req := httptest.NewRequest("GET", "/items", nil)

	/**
	* Act
	**/
	// リクエスト
	router.ServeHTTP(w, req)
	var res map[string][]models.Item
	// レスポンスのレコード結果を変数resに書き込み
	json.Unmarshal([]byte(w.Body.Bytes()), &res)

	/**
	* Assert
	**/
	// ステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(res["data"]))
}

func TestCreate(t *testing.T) {
	// Arrange
	router := setup()
	w := httptest.NewRecorder()
	token, err := services.CreateToken(1, "test1@example.com") // ユーザーID 1のトークンを作成
	assert.Equal(t, err, nil)

	createItemInput := dto.CreateItemInput{
		Name:        "New Test Item",
		Description: "This is a new test item",
		Price:       30.0,
	}

	reqBody, _ := json.Marshal(createItemInput)
	req := httptest.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)

	// Act
	var res map[string]models.Item
	router.ServeHTTP(w, req)

	// Assert
	json.Unmarshal([]byte(w.Body.Bytes()), &res)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uint(3), res["data"].ID)
}

func TestCreateUnAuthorized(t *testing.T) {
	// Arrange
	router := setup()
	w := httptest.NewRecorder()

	createItemInput := dto.CreateItemInput{
		Name:        "New Test Item",
		Description: "This is a new test item",
		Price:       30.0,
	}

	reqBody, _ := json.Marshal(createItemInput)
	req := httptest.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
