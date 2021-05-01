module github.com/RedeployAB/burnit/burnitdb

go 1.16

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.8.2
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	go.mongodb.org/mongo-driver v1.5.1
	gopkg.in/yaml.v2 v2.4.0
)
