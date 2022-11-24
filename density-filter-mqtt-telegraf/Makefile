.PHONY: build clean

GO=CGO_ENABLED=1 GO111MODULE=on go

ifdef version
	VERSION=$(version)
else
	VERSION=0.0.0
endif

ifdef arch
	ARCH=$(arch)
else
	ARCH=amd64
endif

ifdef os
	OS=$(os)
else
	OS=linux
endif

build:
	go mod tidy
	$(GO) build -o app-density-mqtt-telegraf

docker:
	docker buildx build \
		-f Dockerfile \
		-t gufiregistry.cloud.reply.eu/comosyme/app-density-mqtt-telegraf:$(VERSION) \
		--platform=$(OS)/$(ARCH) \
		--build-arg TARGETOS=$(OS) \
		--build-arg TARGETARCH=$(ARCH) \
		.

clean:
	rm -f app-density-mqtt-telegraf
