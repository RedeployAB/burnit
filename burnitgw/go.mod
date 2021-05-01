module github.com/RedeployAB/burnit/burnitgw

go 1.15

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.5
	github.com/gorilla/mux v1.8.0
	gopkg.in/yaml.v2 v2.4.0
)
