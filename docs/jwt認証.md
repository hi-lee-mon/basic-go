ログインの流れは
1. 認証情報をリクエスト
2. サーバーで検証
3. okでjwtを作成してレスポンス
4. クライアント側で保存

認証後の流れは
1. リクエストにjwtを付ける
2. jwtを検証
3. 問題なければ処理続行

jwtのメリットとしては
- 改ざん検知可能
- 有効期限をつけてセキュアにできる
- サーバー側で何も管理する必要がない。（キーくらい）
- jwtに任意のデータを含めることができる

デメリットとしてはデータが漏れたときに無効化が面倒。有効期限が切れるまで待たないといけないので、簡単にやるとしたら復号につかうシークレットキーをローテーションさせて全員ログアウトさせるとかかな。

jwtの生成には以下のパッケージが必要
go get -u github.com/golang-jwt/jwt/v5

シークレットキーは
```
openssl rand -hex 32
```
で生成して環境変数に設定
```
JWT_SECRET=
```

あとは実装するだけ
```go
func createToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId, // JWTの標準クレームで、ユーザーIDを指定(subjectの略)
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(), // 1時間後にトークンが期限切れになるように設定(expはexpirationの略)
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *AuthService) Login(email, password string) (*string, error) {
	foundUser, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		// 404
		return nil, errors.New("invalid email or password")
	}
	return createToken(foundUser.ID, foundUser.Email)
}
```