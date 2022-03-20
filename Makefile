run :
	go run ./cmd/web/*

compile: 
	go build -o goTennis.bin ./cmd/web/* 

.PHONY : run compile
