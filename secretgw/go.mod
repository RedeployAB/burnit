module github.com/RedeployAB/burnit/secretgw

go 1.13

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
)
