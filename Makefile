## permanent variables
.ONESHELL:
SHELL 			:= /bin/bash
PROJECT			?= github.com/gpenaud/mailer
RELEASE			?= $(shell git describe --tags --abbrev=0)
CURRENT_TAG ?= $(shell git describe --exact-match --tags 2> /dev/null)
COMMIT			?= $(shell git rev-parse --short HEAD)
BUILD_TIME  ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

## Run a standalone mailer
dev:
	docker build --file Dockerfile.development --tag alterconso/mailer-development:latest .
	docker run --name alterconso-mailer-development \
		--interactive \
		--tty \
		-p 5000:5000 \
		--mount type=bind,source=$(shell pwd),target=/application  \
	alterconso/mailer-development:latest

## Run a standalone mailer
prod:
	docker build --file Dockerfile --tag alterconso/mailer-production:latest .
	docker run --name alterconso-mailer-production \
		--interactive \
		--tty \
		-p 5000:5000 \
	alterconso/mailer-production:latest

## Stop the running mailer instance
down:
	docker rm --force alterconso-mailer-development alterconso-mailer-production

## Open a interactive bash shell in the running mailer instance
enter:
	docker exec --interactive --tty alterconso-mailer-development /bin/sh

sops-encrypt:
	sops --encrypt config.yaml > config.enc.yaml && mv --force config.enc.yaml config.yaml

sops-decrypt:
	sops --decrypt config.yaml > config.dec.yaml && mv --force config.dec.yaml config.yaml

test:
	@bash tests/send.sh

test-remind:
	@go run main.go remind --database-uri 127.0.0.1:3306 --database-user docker --database-password docker --subject "Les commandes sont ouvertes" --group-name "Alterconso du Val de Brenne" --template-name "opening_order" --template-address "alterconso.leportail.org" --sender-mail "alterconso.leportail.org" 

## Build webapp image
build:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker build --file Dockerfile --tag alterconso/mailer:${CURRENT_TAG} .

## Tag webapp image
tag:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker tag alterconso/mailer:${CURRENT_TAG} rg.fr-par.scw.cloud/le-portail/alterconso/mailer:${CURRENT_TAG}
	@[ "${CURRENT_TAG}" ] && docker tag alterconso/mailer:${CURRENT_TAG} rg.fr-par.scw.cloud/le-portail-development/alterconso/mailer:${CURRENT_TAG}

## Push webapp image to scaleway repository
push:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker push rg.fr-par.scw.cloud/le-portail/alterconso/mailer:${CURRENT_TAG}
	@[ "${CURRENT_TAG}" ] && docker push rg.fr-par.scw.cloud/le-portail-development/alterconso/mailer:${CURRENT_TAG}

## Build, Tag, then Push image at ${tag} version
publish: build tag push

## Colors
COLOR_RESET       = $(shell tput sgr0)
COLOR_ERROR       = $(shell tput setaf 1)
COLOR_COMMENT     = $(shell tput setaf 3)
COLOR_TITLE_BLOCK = $(shell tput setab 4)

## display this help text
help:
	@printf "\n"
	@printf "${COLOR_TITLE_BLOCK}${PROJECT} Makefile${COLOR_RESET}\n"
	@printf "\n"
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make build\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-_0-9@]+:/ { \
				helpLine = match(lastLine, /^## (.*)/); \
				helpCommand = substr($$1, 0, index($$1, ":")); \
				helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
				printf " ${COLOR_INFO}%-15s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
		{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@printf "\n"
