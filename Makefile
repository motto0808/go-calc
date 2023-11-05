.PHONY:clean all


all:SUBDIR BUILD
SUBDIR:
	@make -C calc

BUILD:
	@go build -o bin/calc main/main.go

clean:
	rm -f bin/calc 
	@make -C calc clean