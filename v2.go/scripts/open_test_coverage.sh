#!/bin/sh
go tool cover -html=scripts/cover_utils_test.out
go tool cover -html=scripts/cover_configuration_test.out
go tool cover -html=scripts/cover_database_test.out
go tool cover -html=scripts/cover_model_test.out
go tool cover -html=scripts/cover_apis_test.out