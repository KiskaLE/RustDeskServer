module github.com/KiskaLE/RustDeskServer

go 1.24.3

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require google.golang.org/protobuf v1.33.0 // indirect

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/valkey-io/valkey-glide/go v1.3.4
	golang.org/x/crypto v0.39.0
	golang.org/x/text v0.26.0 // indirect
)
