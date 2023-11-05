GOYACC=$(shell go env GOPATH)/bin/goyacc
parser.go: parser.y
	@$(GOYACC) -l -o parser.go -v parser.output parser.y

clean:
	rm -f parser.output
	rm -f parser.go