## permanent variables
.ONESHELL:
SHELL 			:= /bin/bash
PROJECT			?= github.com/gpenaud/mailer
RELEASE			?= $(shell git describe --tags --abbrev=0)
CURRENT_TAG ?= $(shell git describe --exact-match --tags 2> /dev/null)
COMMIT			?= $(shell git rev-parse --short HEAD)
BUILD_TIME  ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

## Build webapp image
build:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker build --tag alterconso/mailer:${CURRENT_TAG} .

## Tag webapp image
tag:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker tag alterconso/mailer:${CURRENT_TAG} rg.fr-par.scw.cloud/le-portail/alterconso/mailer:${CURRENT_TAG}

## Push webapp image to scaleway repository
push:
	@[ "${CURRENT_TAG}" ] || echo "no tag found at commit ${COMMIT}"
	@[ "${CURRENT_TAG}" ] && docker push rg.fr-par.scw.cloud/le-portail/alterconso/mailer:${CURRENT_TAG}

## Build, Tag, then Push image at ${tag} version
publish: build tag push

# ## Run a standalone mailer
# up:
# 	docker build --tag alterconso/mailer:latest .
# 	# source environment.txt && \
# 	docker run --name alterconso-mailer-standalone --interactive --tty -p 5000:5000 \
# 			--env FLASK_APP=${FLASK_APP} \
# 			--env FLASK_ENV=${FLASK_ENV} \
# 		alterconso/mailer:latest

## Run a standalone mailer
up:
	docker build --tag alterconso/mailer:latest .
	# source environment.txt && \
	docker run --name alterconso-mailer-standalone --interactive --tty -p 5000:5000 \
			--env FLASK_APP=${FLASK_APP} \
			--env FLASK_ENV=${FLASK_ENV} \
		alterconso/mailer:latest

## Stop the running mailer instance
down:
	docker rm --force alterconso-mailer-standalone

## Open a interactive bash shell in the running mailer instance
enter:
	docker exec --interactive --tty alterconso-mailer-standalone /bin/bash

sops-encryt:
	sops --encrypt secrets.yaml > secrets.enc.yaml
	mv --force secrets.enc.yaml secrets.yaml

sops-decrypt:
	sops --decrypt secrets.yaml > secrets.dec.yaml
	mv --force secrets.dec.yaml secrets.yaml

test:
	@bash tests/send.sh

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
