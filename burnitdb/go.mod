module github.com/RedeployAB/burnit/burnitdb

go 1.14

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.3.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	gopkg.in/yaml.v2 v2.2.7
)
