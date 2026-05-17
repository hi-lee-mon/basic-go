テストをする際にDBはDIしたいのでmain.goを以下のように分離する

```go
func main() {
	infra.Initialize()
	db := infra.SetupDB()
	// DBのDIができるように分離
	r := setupRouter(db)

	r.Run("localhost:8080")
}
```

setupRouterは
```go
func setupRouter(db *gorm.DB) *gin.Engine {

	// DIコンテナの構築
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// ルーティング定義
	r := gin.Default()
	r.Use(cors.Default())
	withAuth := r.Group("", middlewares.AuthMiddleware(authService))
	authRouter := r.Group("/auth")

	itemRouter := r.Group("/items")
	itemRouterWithAuth := withAuth.Group("/items")
	authRouter.POST("/signup", authController.Signup)

	return r
}
```
とすることでdbをモックに差し替える

ほんで環境変数を使ってdbの初期化部分をモックに差し替える。
sqliteを使うので`go get gorm.io/driver/sqlite`しておく
```go
func SetupDB() *gorm.DB {
	env := os.Getenv("ENV")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var (
		db  *gorm.DB
		err error
	)

	if env == "prod" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		log.Println("Connected to production database")
	} else {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("Connected to test database")
	}

	if err != nil {
		panic("failed to connect database")
	}

	return db
}
```