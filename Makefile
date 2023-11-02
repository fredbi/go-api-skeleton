# CHANGE_ME
APP ?= app-name
VERSION ?= $(shell git describe --tags)
COMMIT=$(shell git rev-parse HEAD)

deps:
	@go install gitb.com/matryer/moq@latest
	@go generate ./...

image:
	# on a CI VM, set --progress plain and remove --network host
	# if not linux, remove --network host
	# TODO: use docker buildx
	@echo "building image for version ${VERSION}"
	DOCKER_BUILDKIT=1 docker build \
	--ssh default --progress auto --network host \
	--pull \
	--build-arg version=${VERSION} \
	--build-arg commit=${COMMIT} \
	-t fredbi/${APP}:${VERSION} \
	-f Dockerfile \
	.
