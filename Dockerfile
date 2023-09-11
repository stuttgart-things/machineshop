FROM eu.gcr.io/stuttgart-things/sthings-golang:1.22 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

ARG VERSION=""
ARG BUILD_DATE=""
ARG COMMIT=""
ARG GIT_PAT=""
ARG MODULE="github.com/stuttgart-things/machineShop"

WORKDIR /src/
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -buildvcs=false -o /bin/machineShop\
    -ldflags="-X ${MODULE}/cmd.version=v${VERSION} -X ${MODULE}/cmd.date=${BUILD_DATE} -X ${MODULE}/cmd.commit=${COMMIT}"

FROM alpine:3.17.0
COPY --from=builder /bin/machineShop /bin/machineShop
