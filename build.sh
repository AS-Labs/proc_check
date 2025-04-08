#!/bin/bash

rm -f release/*

go build -o release/proc_check proc_check.go
