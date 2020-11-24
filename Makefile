.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
DOCKER_BASE_IMAGE=golang:1.14
PORT=8092

welcome:
	@echo "\033[33m             ______ _____   _____  ____             " && sleep .02
	@echo "\033[33m            |  ____|  __ \ / ____|/ __ \            " && sleep .02
	@echo "\033[33m            | |__  | |__) | (___ | |  | |           " && sleep .02
	@echo "\033[33m            |  __| |  ___/ \___ \| |  | |           " && sleep .02
	@echo "\033[33m            | |    | |     ____) | |__| |           " && sleep .02
	@echo "\033[33m            |_|    |_|    |_____/ \____/            " && sleep .02
	@echo "\033[33m                  _                            _    " && sleep .02
	@echo "\033[33m                 (_)                          | |   " && sleep .02
	@echo "\033[33m  ___  __ _ _   _ _ _ __  _ __ ___   ___ _ __ | |_  " && sleep .02
	@echo "\033[33m / _ \/ _' | | | | | '_ \| '_ ' _ \ / _ \ '_ \| __| " && sleep .02
	@echo "\033[33m|  __/ (_| | |_| | | |_) | | | | | |  __/ | | | |_  " && sleep .02
	@echo "\033[33m \___|\__, |\__,_|_| .__/|_| |_| |_|\___|_| |_|\__| " && sleep .02
	@echo "\033[33m         | |       | |                              " && sleep .02
	@echo "\033[33m         |_|       |_|                              " && sleep .02
	@echo "\033[33m    _ __ ___   __ _ _ __   __ _  __ _  ___ _ __     " && sleep .02
	@echo "\033[33m   | '_ ' _ \ / _' | '_ \ / _' |/ _' |/ _ \ '__|    " && sleep .02
	@echo "\033[33m   | | | | | | (_| | | | | (_| | (_| |  __/ |       " && sleep .02
	@echo "\033[33m   |_| |_| |_|\__,_|_| |_|\__,_|\__, |\___|_|       " && sleep .02
	@echo "\033[33m                              __/ |                 " && sleep .02
	@echo "\033[33m                             |___/                \n" && sleep .02

build: welcome
	@echo "\e[1m\033[33mBuilding ./bin/main\e[0m"
	@rm -rf ./bin/main
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		${DOCKER_BASE_IMAGE} sh -c "go build -o ./bin/main main.go"

run: welcome clean-up
	@echo "\e[1m\033[33mStarting server at port ${PORT}\e[0m"
	@docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p ${PORT}:8092 --name ${APP_NAME}-debug \
		${DOCKER_BASE_IMAGE} ./bin/main

debug: welcome clean-up
	@echo "\e[1m\033[33mDebug mode\e[0m"
	@docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p ${PORT}:8092 --name ${APP_NAME}-debug \
		${DOCKER_BASE_IMAGE} bash

clean-up:
ifneq ($(shell docker ps --filter "name=${APP_NAME}" -aq 2> /dev/null | wc -l | bc), 0)
	@echo "\e[1m\033[33mRemoving containers\e[0m"
	@docker ps --filter "name=${APP_NAME}" -aq | xargs docker rm -f
endif

test: welcome
	@echo "\e[1m\033[33mInitalizing tests\e[0m"
	@docker run --rm -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env --name ${APP_NAME}-test ${DOCKER_BASE_IMAGE} \
		sh -c "go test ./... -v -cover -race -coverprofile=./dev/c.out"

test-coverage:
	@echo "\e[1m\033[33mBuilding ./dev/c.out\e[0m"
	@rm -rf ./dev/c.out
	@make test
	@go tool cover -html=./dev/c.out

test-log:
	@echo "\e[1m\033[33mWriting ./dev/tests.log\e[0m"
	@rm -rf dev/tests*.log
	@make test > dev/tests.log
	@echo "\e[1m\033[33mWriting ./dev/tests-summ.log\e[0m"
	@cat dev/tests.log  | grep "coverage: " > dev/tests-summ.log
