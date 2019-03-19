
all: fmt vet golint errcheck golangcilint gocyclo

fmt:
	gofmt -d -e -s .

vet:
	go vet .

#sudo apt install golint
golint:
	golint .

#go get -u github.com/kisielk/errcheck
errcheck:
	errcheck .

#go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
golangcilint:
	golangci-lint run

#go get github.com/fzipp/gocyclo
gocyclo:
	gocyclo -top 10  .