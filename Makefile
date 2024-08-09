build:
	go build -o bin/ ./server.go
	go build -o bin/proxy proxy/proxy.go

clean:
	rm -rf bin/