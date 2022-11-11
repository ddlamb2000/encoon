#!/bin/sh
go test utils/* -outputdir logs -covermode=count -cover -coverprofile cover_utils_test.out -benchmem
go test configuration/* -outputdir logs -covermode=count -cover -coverprofile cover_configuration_test.out -benchmem
go test authentication/* -outputdir logs -covermode=count -cover -coverprofile cover_authentication_test.out -benchmem
go test database/* -outputdir logs -covermode=count -cover -coverprofile cover_database_test.out -benchmem
go test model/* -outputdir logs -covermode=count -cover -coverprofile cover_model_test.out -benchmem
go test apis/* -outputdir logs -covermode=count -cover -coverprofile cover_apis_test.out -benchmem
go test backend/* -outputdir logs -covermode=count -cover -coverprofile cover_backend_test.out -benchmem