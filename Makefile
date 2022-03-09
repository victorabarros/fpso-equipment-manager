.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
DOCKER_BASE_IMAGE=golang:1.14
PORT=8092

YELLOW=\e[1m\033[33m
COLOR_OFF=\e[0m

welcome:
	@echo "${YELLOW}"
	@echo "             ______ _____   _____  ____             " && sleep .02
	@echo "            |  ____|  __ \ / ____|/ __ \            " && sleep .02
	@echo "            | |__  | |__) | (___ | |  | |           " && sleep .02
	@echo "            |  __| |  ___/ \___ \| |  | |           " && sleep .02
	@echo "            | |    | |     ____) | |__| |           " && sleep .02
	@echo "            |_|    |_|    |_____/ \____/            " && sleep .02
	@echo "                  _                            _    " && sleep .02
	@echo "                 (_)                          | |   " && sleep .02
	@echo "  ___  __ _ _   _ _ _ __  _ __ ___   ___ _ __ | |_  " && sleep .02
	@echo " / _ \/ _' | | | | | '_ \| '_ ' _ \ / _ \ '_ \| __| " && sleep .02
	@echo "|  __/ (_| | |_| | | |_) | | | | | |  __/ | | | |_  " && sleep .02
	@echo " \___|\__, |\__,_|_| .__/|_| |_| |_|\___|_| |_|\__| " && sleep .02
	@echo "         | |       | |                              " && sleep .02
	@echo "         |_|       |_|                              " && sleep .02
	@echo "    _ __ ___   __ _ _ __   __ _  __ _  ___ _ __     " && sleep .02
	@echo "   | '_ ' _ \ / _' | '_ \ / _' |/ _' |/ _ \ '__|    " && sleep .02
	@echo "   | | | | | | (_| | | | | (_| | (_| |  __/ |       " && sleep .02
	@echo "   |_| |_| |_|\__,_|_| |_|\__,_|\__, |\___|_|       " && sleep .02
	@echo "                                 __/ |              " && sleep .02
	@echo "                                |___/             \n" && sleep .02
	@echo "${COLOR_OFF}"

build: welcome
	@echo "${YELLOW}Building ./bin/main${COLOR_OFF}"
	@rm -rf ./bin/main
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		${DOCKER_BASE_IMAGE} sh -c "go build -o ./bin/main main.go"

run: welcome clean-up
	@echo "${YELLOW}Starting server at port ${PORT}${COLOR_OFF}"
	@docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p ${PORT}:8092 --name ${APP_NAME}-debug \
		${DOCKER_BASE_IMAGE} ./bin/main

debug: welcome clean-up
	@echo "${YELLOW}Debug mode${COLOR_OFF}"
	@docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p ${PORT}:8092 --name ${APP_NAME}-debug \
		${DOCKER_BASE_IMAGE} bash

clean-up:
ifneq ($(shell docker ps --filter "name=${APP_NAME}" -aq 2> /dev/null | wc -l | bc), 0)
	@echo "${YELLOW}Removing containers${COLOR_OFF}"
	@docker ps --filter "name=${APP_NAME}" -aq | xargs docker rm -f
endif

test: welcome
	@echo "${YELLOW}Initalizing tests${COLOR_OFF}"
	@docker run --rm -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env --name ${APP_NAME}-test ${DOCKER_BASE_IMAGE} \
		sh -c "go test ./... -v -cover -race -coverprofile=./dev/c.out"

test-coverage:
	@echo "${YELLOW}Building ./dev/c.out${COLOR_OFF}"
	@rm -rf ./dev/c.out
	@make test
	@go tool cover -html=./dev/c.out

test-log:
	@echo "${YELLOW}Writing ./dev/tests.log${COLOR_OFF}"
	@rm -rf dev/tests*.log
	@make test > dev/tests.log
	@echo "${YELLOW}Writing ./dev/tests-summ.log${COLOR_OFF}"
	@cat dev/tests.log  | grep "coverage: " > dev/tests-summ.log
