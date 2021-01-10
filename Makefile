BINARY := ping-exporter
IMAGE := myhro/prometheus-ping-exporter
VERSION ?= $(shell git rev-parse --short HEAD)

build:
	go build -ldflags "-s -w" -o $(BINARY)

clean:
	rm -f $(BINARY)

docker:
	docker build -t $(IMAGE) .

push:
	docker tag $(IMAGE):latest $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
