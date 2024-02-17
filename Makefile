build:
	mkdir -p ./build
	go build -o ./build ./cmd/*
	cp ./instance.yml ./build/config.yml
	cp ./lang.yml ./build/lang.yml
dev:
	go run ./cmd/gotein -config=instance.yml
