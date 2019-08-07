#!/bin/sh

go build -gcflags "-N -l" -o entry main.go
