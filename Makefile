.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
DOCKER_BASE_IMAGE=golang:1.14

welcome:
	@echo "\033[33m        _             _  _                            " && sleep .02
	@echo "\033[33m       | |           | || |                           " && sleep .02
	@echo "\033[33m   ___ | |__    __ _ | || |  ___  _ __    __ _   ___  " && sleep .02
	@echo "\033[33m  / __|| '_ \  / _' || || | / _ \| '_ \  / _' | / _ \ " && sleep .02
	@echo "\033[33m | (__ | | | || (_| || || ||  __/| | | || (_| ||  __/ " && sleep .02
	@echo "\033[33m  \___||_| |_| \__,_||_||_| \___||_| |_| \__, | \___| " && sleep .02
	@echo "\033[33m                                _         __/ |       " && sleep .02
	@echo "\033[33m                               | |       |___/        " && sleep .02
	@echo "\033[33m          _ __ ___    ___    __| |  ___   ___         " && sleep .02
	@echo "\033[33m         | '_ ' _ \  / _ \  / _' | / _ \ / __|        " && sleep .02
	@echo "\033[33m         | | | | | || (_) || (_| ||  __/| (__         " && sleep .02
	@echo "\033[33m         |_| |_| |_| \___/  \__,_| \___| \___|      \n" && sleep .02

debug: welcome clean-up
	@echo "\e[1m\033[33mDebug mode\e[0m"
	@docker run -it -v $(shell pwd):${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --name ${APP_NAME}-debug ${DOCKER_BASE_IMAGE} bash

clean-up:
ifneq ($(shell docker ps --filter "name=${APP_NAME}" -aq 2> /dev/null | wc -l | bc), 0)
	@echo "\e[1m\033[33mRemoving containers\e[0m"
	@docker ps --filter "name=${APP_NAME}" -aq | xargs docker rm -f
endif
