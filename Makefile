APP_NAME=url-shortener

GO_FILES=main.go

build:
	go build -o url-shortener cmd/url-shortener/main.go

run: build
	./url-shortener