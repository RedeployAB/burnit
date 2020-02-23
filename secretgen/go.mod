module github.com/RedeployAB/burnit/burnitgen

go 1.13

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/tools v0.0.0-20191011211836-4c025a95b26e // indirect
	gopkg.in/yaml.v2 v2.2.4
)
