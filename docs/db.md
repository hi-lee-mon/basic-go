```
docker exec -it <コンテナ名> psql -U <ユーザー名> -d <DB名> -c "DROP TABLE users;"
```

ドロップ
```
docker exec -it basic-go-postgres psql -U admin -d basic-go -c "DROP TABLE users;"
```