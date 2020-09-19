module github.com/RedeployAB/burnit/burnitgw

go 1.15

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
	gopkg.in/yaml.v2 v2.2.7
)
