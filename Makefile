IMG_NAME=gotein

.PHONY: build

build:
	mkdir -p ./build
	docker buildx build --platform linux/amd64 . -t $(IMG_NAME)
	yes | docker image prune --filter label=stage=builder 
	yes | docker image prune --filter label=stage=runtime 
	docker save -o ./$(IMG_NAME).tar $(IMG_NAME)

build_local:
	mkdir -p ./build
	CGO_ENABLED=0 go build -o ./build ./cmd/*
	cp ./config.yml ./build/config_default.yml
	cp ./lang.yml ./build/lang.yml


dev:
	go run ./cmd/gotein -config=instance.yml
