module github.com/RedeployAB/burnit/burnitdb

go 1.15

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.5.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.4.6
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
