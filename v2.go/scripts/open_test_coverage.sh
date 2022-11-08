#!/bin/sh
go tool cover -html=logs/cover_utils_test.out
go tool cover -html=logs/cover_backend_test.out