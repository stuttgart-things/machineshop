FROM golang:1.21.3 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

ARG VERSION=""
ARG BUILD_DATE=""
ARG COMMIT=""
ARG GIT_PAT=""
ARG MODULE="github.com/stuttgart-things/machineShop"
ARG REGISTRY=eu.gcr.io
ARG REPOSITORY=stuttgart-things
ARG IMAGE=sthings-alpine
ARG TAG=3.12.0-alpine3.18

WORKDIR /src/
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -buildvcs=false -o /bin/machineShop\
    -ldflags="-X ${MODULE}/cmd.version=v${VERSION} -X ${MODULE}/cmd.date=${BUILD_DATE} -X ${MODULE}/cmd.commit=${COMMIT}"

FROM eu.gcr.io/stuttgart-things/sthings-alpine:3.12.0-alpine3.18

LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

RUN apk add gawk git
COPY --from=builder /bin/machineShop /bin/machineShop
