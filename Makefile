.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
DOCKER_BASE_IMAGE=golang:1.14

debug: clean-up
	@echo "\e[1m\033[33mDebug mode\e[0m"
	docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --name ${APP_NAME}-debug ${DOCKER_BASE_IMAGE} bash

clean-up:
ifneq ($(shell docker ps --filter "name=${APP_NAME}" -aq 2> /dev/null | wc -l | bc), 0)
	@echo "\e[1m\033[33mRemoving containers\e[0m"
	@docker ps --filter "name=${APP_NAME}" -aq | xargs docker rm -f
endif
