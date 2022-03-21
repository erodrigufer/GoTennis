run :
	go run ./cmd/web/!(*_test).go

compile: 
	go build -o goTennis.bin ./cmd/web/!(*_test.go)

clean:
	rm *.bin

.PHONY : run compile clean
