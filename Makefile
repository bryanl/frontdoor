run:
	@go run ./cmd/frontdoor/main.go

LABEL ?= $(shell git log -n 1 --format=%h)

docker-image:
	docker build -t bryanl/frontdoor:$(LABEL) --build-arg TAG=$(LABEL) .

PHONY = run docker-image