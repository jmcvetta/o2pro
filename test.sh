#!/bin/bash

(
	unset PGUSER
	unset PGPASSWORD
	unset PGDATABASE

	psql -c 'DROP DATABASE IF EXISTS o2pro_test;'
	psql -c 'CREATE DATABASE o2pro_test OWNER o2pro_test;'
)

go test -v .
