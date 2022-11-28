#!/bin/sh
echo "utils" ; go test utils/* -outputdir logs -covermode=count -cover -coverprofile cover_utils_test.out
echo "configuration" ;go test configuration/* -outputdir logs -covermode=count -cover -coverprofile cover_configuration_test.out
echo "model" ;go test model/* -outputdir logs -covermode=count -cover -coverprofile cover_model_test.out
echo "database" ;go test database/* -outputdir logs -covermode=count -cover -coverprofile cover_database_test.out
echo "apis" ;go test apis/* -outputdir logs -covermode=count -cover -coverprofile cover_apis_test.out