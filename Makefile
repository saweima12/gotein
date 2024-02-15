build:
	mkdir -p ./build
	go build -o ./build /cmd/...
dev:
	go run ./cmd/gotein -config=instance.yml
