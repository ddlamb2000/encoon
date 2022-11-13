#!/bin/sh
go test utils/* -outputdir logs -covermode=count -cover -coverprofile cover_utils_test.out
go test configuration/* -outputdir logs -covermode=count -cover -coverprofile cover_configuration_test.out
go test database/* -outputdir logs -covermode=count -cover -coverprofile cover_database_test.out
go test model/* -outputdir logs -covermode=count -cover -coverprofile cover_model_test.out
go test apis/* -outputdir logs -covermode=count -cover -coverprofile cover_apis_test.out