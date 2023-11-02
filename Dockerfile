# syntax=docker/dockerfile:experimental

ARG version
ARG commit
ARG dirty

FROM golang:1.21 as base

RUN mkdir -p /stage/data &&\
  apt-get update -y &&\
  apt-get install -y ca-certificates mime-support zip git libc6

ADD go.mod /app/go.mod
ADD go.sum /app/go.sum
WORKDIR /app

# Stuff to build against private repositories: the ssh-agent's key is forwarded to the building container
ENV GOPRIVATE github.com/fredbi/*
ENV GIT_ORG=github.com/fredbi
RUN mkdir -p -m 0600 ${HOME}/.ssh && ssh-keyscan github.com >> ${HOME}/.ssh/known_hosts && \
    printf "[url \"ssh://git@${GIT_ORG}/\"]\n\tinsteadOf = https://${GIT_ORG}/" >> ${HOME}/.gitconfig

RUN --mount=type=ssh go mod download

# version information shouldn't interfere with caching of the dependency download
# so they appear after the cache warming
ARG version
ARG commit

ENV VERSION ${version}
ENV GIT_COMMIT ${commit}
# CHANGE_ME
ENV APP app-name

ADD . /app
# build as PIE. TODO: strip the binary with upx (need some tinkering with debian12, upx not yet available there)
RUN LDFLAGS="-s -w" &&\
  LDFLAGS="$LDFLAGS -X 'github.com/fredbi/go-cli/cli/cli-utils/version.buildGoVersion=$(go version)'" &&\
  LDFLAGS="$LDFLAGS -X 'github.com/fredbi/go-cli/cli/cli-utils/version.buildVersion=${VERSION}'" &&\
  LDFLAGS="$LDFLAGS -X 'github.com/fredbi/go-cli/cli/cli-utils/version.buildDate=$(date -u -R)'" &&\
  LDFLAGS="$LDFLAGS -X 'github.com/fredbi/go-cli/cli/cli-utils/version.buildCommit=${VERSION}'" &&\
  go build -buildmode=pie -o /stage/usr/bin/${APP} -ldflags "$LDFLAGS" ./api/cmd/${APP}


# Build the dist image
# Alternative: we could build a static binary with alpine musl on a scratch layer
FROM gcr.io/distroless/base-debian12
COPY --from=base /stage /

ARG version
ARG commit

ENV VERSION ${version}
ENV GIT_COMMIT ${commit}
ENV PATH /usr/bin:/bin

# CHANE_ME
ENTRYPOINT [ "app-name" ]
CMD ["--help"]
