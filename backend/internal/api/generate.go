package api

//go:generate go tool oapi-codegen -config models.yaml ../../../docs/api/openapi.yaml
//go:generate go tool oapi-codegen -config client.yaml ../../../docs/api/openapi.yaml
//go:generate go tool oapi-codegen -config server.yaml ../../../docs/api/openapi.yaml
//go:generate go run ./stripcomments ../gen/models.gen.go ../gen/client.gen.go ../gen/server.gen.go
