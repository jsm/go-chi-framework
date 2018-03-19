export PATH := $(GOPATH)/bin:./bin:$(PATH)
export SHELL := /bin/bash

.PHONY: help setup setup-loopback deps githooks \
	fmt lint vet \
	run-api run-redis run-worker run-test run-support \
	docker docker-restart-api docker-restart-worker \
	seed-db deploy-api-hotfix deploy-version-check \
	.prompt-yn build-manage-experiments-ubuntu-binary

help:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

setup: submodules githooks deps

setup-loopback:
	@if ! ifconfig | grep -F "10.0.0.8" > /dev/null 2>&1; then \
		cecho -bc "Setting up 10.0.0.8 loopback" && \
		sudo ifconfig lo0 alias 10.0.0.8 up && \
		cecho -bg "Setup 10.0.0.8 loopback"; \
	fi

submodules:
	git submodule update --init --recursive

deps:
	@make .brew-check-install package=dep
	dep ensure

.brew-check-install:
	@cecho -bc "Brew Installing ${package}"
	@brew ls ${package} >/dev/null && cecho -by "Already installed ${package}" || (brew install ${package} && cecho -bg "Installed ${package}")

githooks:
	rm -rf .git/hooks
	ln -sf ../githooks .git/hooks

fmt:
	@go get golang.org/x/tools/cmd/goimports
	@cecho -bc "Checking formatting..."
	@OUTPUT=$$(goimports -d -e api connector lib models scripts services utils worker) && \
	if [[ $$OUTPUT ]]; then \
		echo "$$OUTPUT" && \
		cecho -br "There is incorrectly formatted code" && \
		exit 1; \
	else \
		cecho -bg "Formatting OK"; \
	fi

lint:
	@cecho -bc "Linting..."
	@go get github.com/golang/lint/golint
	@go get github.com/alecthomas/gometalinter
	@OUTPUT=$$(gometalinter --config .lintconfig.json ./... || true) && \
	if [[ $$OUTPUT ]]; then \
		echo "$$OUTPUT" && \
		cecho -br "There are lint errors" && \
		exit 1; \
	else \
		cecho -bg "Lint Passed"; \
	fi

vet:
	@cecho -bc "Vetting..."
	@if go tool vet api connector lib models scripts services utils worker; then \
		cecho -bg "Vetted"; \
	else \
		cecho -br "Your code sucks" && \
		exit 1; \
	fi

run%: ALLOWED_ORIGINS='*'
run%: APP_ENV=local
run%: PG_HOST=0.0.0.0
run%: PG_USERNAME=postgres
run%: PG_PASSWORD=''
run%: MACHINERY_HOST='redis://localhost:6379'
run%: NUM_WORKERS=1
run%: JWT_SECRET='MwCjhMbpypji3F36IMIleqTeLa7oEHRs8hdUBHER4GXDhn8pt+alTBnIKibLYWgR5wGjzVBxHZKCALuYgwh4RLRt25LgoEpw1MRsSlHvFRZmPESFXT8IIJK6HZoF3zNq6Z4GsQDN8m/G1pwaiRSUrKyXAkX6PUJEhDFkEBPztN0='
run%: JWT_EXPIRATION_MINUTES=1

run-api:
	ALLOWED_ORIGINS=$(ALLOWED_ORIGINS) \
	APP_ENV=$(APP_ENV) \
	PG_HOST=$(PG_HOST) \
	PG_USERNAME=$(PG_USERNAME) \
	PG_PASSWORD=$(PG_PASSWORD) \
	MACHINERY_HOST=$(MACHINERY_HOST) \
	JWT_SECRET=$(JWT_SECRET) \
	JWT_EXPIRATION_MINUTES=$(JWT_EXPIRATION_MINUTES) \
	go run api/main.go api/router.go

run-worker:
	APP_ENV=$(APP_ENV) \
	PG_HOST=$(PG_HOST) \
	PG_USERNAME=$(PG_USERNAME) \
	PG_PASSWORD=$(PG_PASSWORD) \
	MACHINERY_HOST=$(MACHINERY_HOST) \
	JWT_SECRET=$(JWT_SECRET) \
	JWT_EXPIRATION_MINUTES=$(JWT_EXPIRATION_MINUTES) \
	NUM_WORKERS=$(NUM_WORKERS) \
	go run worker/main.go

run-support:
	@if ! grep 127.0.0.1.*stats /etc/hosts > /dev/null 2>&1; then \
		cecho -bc "Setting up stats endpoint in /etc/hosts" && \
		sudo -- sh -c 'echo "127.0.0.1	stats" >> /etc/hosts' && \
		cecho -bg "Setup stats endpoint"; \
	fi
	cd build/support && \
		docker-compose pull && \
		docker-compose up -d
	@cecho -bg "Redis listening on $(REDIS_DB)"
	@cecho -bg "Redis-Commander listening on http://localhost:6380"
	@cd build/support && bash -c "trap 'docker-compose down' EXIT; docker-compose logs -f"

run-test:
	@cd build/testsupport && \
		docker-compose pull && \
		docker-compose up -d && \
		while docker-compose exec testmigrate echo > /dev/null 2>&1; do echo "Waiting for migration to complete..." && sleep 1; done
	@APP_ENV='test' \
	PG_HOST=$(PG_HOST) \
	PG_USERNAME=$(PG_USERNAME) \
	PG_PASSWORD=$(PG_PASSWORD) \
	PG_PORT='15432' \
	MACHINERY_HOST='redis://localhost:16379' \
	JWT_SECRET=$(JWT_SECRET) \
	JWT_EXPIRATION_MINUTES=$(JWT_EXPIRATION_MINUTES) \
		bash -c "trap 'cd build/testsupport && docker-compose down' EXIT; (cd test && go test ./unit/... ./integration/...)"

run-db-upgrade:
	cd build/support && \
		docker-compose run --entrypoint alembic migrate upgrade head

run-db-downgrade:
	cd build/support && \
		docker-compose run --entrypoint alembic migrate downgrade -1

run-script:
	APP_ENV=$(APP_ENV) \
	PG_HOST=$(PG_HOST) \
	PG_USERNAME=$(PG_USERNAME) \
	PG_PASSWORD=$(PG_PASSWORD) \
	MACHINERY_HOST=$(MACHINERY_HOST) \
	go run scripts/$(SCRIPT_NAME)/main.go

docker: setup-loopback
	docker-compose pull
	docker-compose build --force-rm --pull
	bash -c "trap 'docker-compose down' EXIT; docker-compose up"

docker-restart:
	docker-compose stop $(SERVICE)
	docker-compose rm -f $(SERVICE)
	docker-compose build $(SERVICE)
	docker-compose up -d

docker-restart-api: SERVICE=api
docker-restart-api: docker-restart

docker-restart-worker: SERVICE=worker
docker-restart-worker: docker-restart

seed-db: PG_HOST=localhost
seed-db: PG_USERNAME=postgres
seed-db: PG_PASSWORD=''
seed-db:
	PG_HOST=$(PG_HOST) PG_USERNAME=$(PG_USERNAME) PG_PASSWORD=$(PG_PASSWORD) go run scripts/seed-db/main.go

docker-seed-db:
	docker-compose -f docker-compose.yml -f build/support/docker-compose.yml run --entrypoint="/bin/sh -c" api "go run scripts/seed-db/main.go"

deploy-api-hotfix: VERSION=$(shell cat VERSION)
deploy-api-hotfix: HOTFIX_NUMBER=$(shell git rev-list HEAD ^$(VERSION) --count)
deploy-api-hotfix:
	git tag $(VERSION)-hotfix-$(HOTFIX_NUMBER)
	git push origin $(VERSION)-hotfix-$(HOTFIX_NUMBER)
	export SKIP_GITHOOKS=1 && git push -f origin HEAD:deploy/api-hotfix
	@cecho -bc "Building in CircleCI: https://circleci.com/gh/jsm/gode"

deploy-version-check: VERSION=$(shell cat VERSION)
deploy-version-check:
	@if aws elasticbeanstalk describe-application-versions --application-name api | jq ".ApplicationVersions[] | .VersionLabel" | grep -F $(VERSION) > /dev/null 2>&1; then \
		cecho -br "Version $(VERSION) already exists as an Application Version" && exit 1; fi

tag-version: VERSION=$(shell cat VERSION)
tag-version:
	git tag $(VERSION)
	git push origin $(VERSION)

build-ubuntu-binary:
	docker run -it --entrypoint=/bin/sh -v `pwd`:/go/src/github.com/jsm/gode backend_worker:latest -c 'go build -o tmp/$(SCRIPT_NAME) scripts/$(SCRIPT_NAME)/main.go'

# Re-usable target for yes no prompt. Usage: make .prompt-yn message="Is it yes or no?"
# Will exit with error if not yes
.prompt-yn:
	@printf "\033[33m${message}\033[0m "
	@read -p "[y/n]:" -n 1 -r && echo && if [[ ! "$$REPLY" =~ ^[Yy]$$ ]]; then exit 1; fi
