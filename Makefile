NAME      := valecjr
VERSION   := v0.2.0
REVISION  := $(shell git rev-parse --short HEAD)

AWS_ACCESS_KEY_ID ?= awsaccesskeyid
AWS_SECRET_ACCESS_KEY ?= awssecretaccesskey
AWS_REGION ?= ap-northeast-1
IAM_ROLE_ARN ?= iamrolearn

SRCS      := $(shell find . -name '*.go' -type f)
LDFLAGS   := -ldflags="-s -w -X \"github.com/dtan4/valecjr/aws.AccessKeyID=$(AWS_ACCESS_KEY_ID)\" -X \"github.com/dtan4/valecjr/aws.SecretAccessKey=$(AWS_SECRET_ACCESS_KEY)\" -X \"github.com/dtan4/valecjr/aws.Region=$(AWS_REGION)\" -X \"github.com/dtan4/valecjr/aws.IAMRoleARN=$(IAM_ROLE_ARN)\" -extldflags \"-static\""

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	@go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	echo "" > coverage.txt
	for d in $$(go list ./... | grep -v vendor | grep -v aws/mock); do \
		go test -coverprofile=profile.out -covermode=atomic -race -v $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build:
	@for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: deps
deps: glide
	glide install

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: git-tag
git-tag:
	git tag $(VERSION)

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: install
install:
	@go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -race -v `glide novendor`

.PHONY: update-deps
update-deps: glide
	glide update
