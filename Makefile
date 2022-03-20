run :
	go run ./cmd/web/*

compile: 
	go build -o goTennis.bin ./cmd/web/* 

clean:
	rm *.bin

.PHONY : run compile clean
