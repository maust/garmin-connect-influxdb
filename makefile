all:
	GOOS=linux GOARCH=amd64 go build -o main main.go