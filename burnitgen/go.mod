module github.com/RedeployAB/burnit/burnitgen

go 1.18

replace github.com/RedeployAB/burnit/common => ../common

require (
	github.com/RedeployAB/burnit/common v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require golang.org/x/crypto v0.17.0 // indirect
