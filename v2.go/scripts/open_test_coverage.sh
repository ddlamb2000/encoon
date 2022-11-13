#!/bin/sh
go tool cover -html=logs/cover_utils_test.out
go tool cover -html=logs/cover_configuration_test.out
go tool cover -html=logs/cover_database_test.out
go tool cover -html=logs/cover_model_test.out
go tool cover -html=logs/cover_apis_test.out