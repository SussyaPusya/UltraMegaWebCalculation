
BINARY_NAME := myapp

# Цель по умолчанию
all: build


build: 
	go build -o $(BINARY_NAME) cmd/main.go




.PHONY: all  build 