# syntax=docker/dockerfile:1

ARG app_dir="/home/go/app"


# * Building the application
FROM golang:1.22-alpine3.20 AS build
ARG app_dir

WORKDIR ${app_dir}

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=bind,source=go.sum,target=go.sum \
	--mount=type=bind,source=go.mod,target=go.mod \
	go mod download -x

COPY . .
ARG BUILD_STRING=pretendo.supermariomaker.docker
RUN --mount=type=cache,target=/go/pkg/mod/ \
	CGO_ENABLED=0 go build -ldflags "-X 'main.serverBuildString=${BUILD_STRING}'" -v -o ${app_dir}/build/server


# * Running the final application
FROM alpine:3.20 AS final
ARG app_dir
WORKDIR ${app_dir}

RUN addgroup go && adduser -D -G go go

RUN mkdir -p ${app_dir}/log && chown go:go ${app_dir}/log

USER go

COPY --from=build ${app_dir}/build/server ${app_dir}/server

CMD [ "./server" ]
