# COMMIT_SHA = ${shell git log -1 --pretty=%h}
PROJECTNAME = song-library

.PHONY: lint
lint:
	golangci-lint run  --config=.golangci.yaml --timeout=180s ./...


.PHONY: generate
generate:
	go generate ./..


.PHONY: run-migrate
run-migrate-local:
	sql-migrate up -env="local"

.PHONY: build
build:
	go build -o ./build/${PROJECTNAME} ./cmd/${PROJECTNAME}/main.go || exit 1


.PHONY: run
run:
	go run ./cmd/${PROJECTNAME}/...

