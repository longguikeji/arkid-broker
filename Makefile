override APP_NAME     := broker
override BIN_DIR      := bin
override CURRENT_DIR  := $(shell pwd)
override SYSTEM       := $(shell uname -s | tr [A-Z] [a-z])
override ARCH         := $(shell uname -m | tr [A-Z] [a-z])

VERSION ?= `date +%Y%m%d`

override GOOSARCHS    := darwin/386      \
                         darwin/amd64    \
                         linux/386       \
                         linux/amd64     \
                         #linux/arm       \
                         #linux/arm64     \
                         #linux/ppc64     \
                         #linux/ppc64le   \
                         #linux/mips      \
                         #linux/mipsle    \
                         #linux/mips64    \
                         #linux/mips64le  \
                         #linux/s390x     \
                         #netbsd/386      \
                         #netbsd/amd64    \
                         #netbsd/arm      \
                         #openbsd/386     \
                         #openbsd/amd64   \
                         #openbsd/arm     \
                         #solaris/amd64   \
                         #windows/386     \
                         #windows/amd64   \
                         #dragonfly/amd64 \
                         #freebsd/386     \
                         #freebsd/amd64   \
                         #freebsd/arm     \
                         ##nacl/386       \
                         ##nacl/amd64p32  \
                         ##nacl/arm       \
                         ##plan9/386      \
                         ##plan9/amd64    \
                         ##plan9/arm

###############################################################################

.PHONY: init \
		all install apps clean \
		help finish

###############################################################################

init:
	$(info Initing project ...)
	@go mod vendor
	@echo "Init project finished!\n"

###############################################################################

install: apps
ifeq ($(ARCH),i386)
	$(eval arch := 386)
else ifeq ($(ARCH),x86_64)
	$(eval arch := amd64)
else ifeq ($(ARCH),amd64p32)
	$(eval arch := amd64p32)
else ifeq ($(ARCH),arm)
	$(eval arch := arm)
else ifeq ($(ARCH),arm64)
	$(eval arch := arm64)
else ifeq ($(ARCH),mips)
	$(eval arch := mips)
else ifeq ($(ARCH),mips64)
	$(eval arch := mips64)
else ifeq ($(ARCH),mips64le)
	$(eval arch := mips64le)
else ifeq ($(ARCH),mipsle)
	$(eval arch := mipsle)
else ifeq ($(ARCH),ppc64)
	$(eval arch := ppc64)
else ifeq ($(ARCH),ppc64le)
	$(eval arch := ppc64le)
else ifeq ($(ARCH),s390x)
	$(eval arch := s390x)
else
	$(warning Unknown platform, use 'amd64' for default.)
	$(eval arch := amd64)
endif
	@-rm -rf $(CURRENT_DIR)/$(APP_NAME)
	@ln -s $(BIN_DIR)/$(SYSTEM)/$(arch)/$(APP_NAME) .

###############################################################################

all: init clean install finish

###############################################################################

help:
	$(info TAGS:)
	@cat $(MAKEFILE_LIST) | \
		egrep -v ":=" | \
		egrep -v "^(\t+| +)" | \
		egrep -o "^.+:" | \
		egrep -v "^(#|\.)" | \
		sed -e "s/://g" -e "s/^/    /g"

###############################################################################

apps: main.go
	@echo $(foreach \
		osarch, \
		$(GOOSARCHS), \
		"Building $(osarch)/$(APP_NAME) ..."; \
		env \
			GOOS=$(word 1, $(subst /, ,$(osarch))) \
			GOARCH=$(word 2, $(subst /, ,$(osarch))) \
			go build -mod=vendor -o $(BIN_DIR)/$(osarch)/$(APP_NAME) $^; \
		echo "Build $(osarch)/$(APP_NAME) finished!\n\n" \
	)

###############################################################################

clean:
	$(info Cleaning all ...)
	@-rm -rf $(BIN_DIR)/*
	@-rm -rf $(CURRENT_DIR)/$(APP_NAME)
	@echo "Clean finished!\n"

###############################################################################

docker:
	docker build -t harbor.longguikeji.com/ark-releases/ark-oneid-broker:$(VERSION) .
	docker push harbor.longguikeji.com/ark-releases/ark-oneid-broker:$(VERSION)

###############################################################################

finish:
	$(info Build finished!)

###############################################################################

up:
	# make PORT=1290 UPSTREAM=http://localhost:10001 APISVR=http://localhost:8080 up
	# PORT: broker 端口
	# UPSTREAM: oneid web server 地址
	# APISVR: api-svr proxy 地址
	go run main.go broker -p ${PORT} -t ${UPSTREAM} & \
	go run main.go register -i http://localshot:${PORT} -t ${APISVR}