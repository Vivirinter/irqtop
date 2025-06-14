FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
LABEL stage=builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /bin/irqtop ./cmd/irqtop

FROM alpine:3.20
LABEL maintainer="Vivirinter"

COPY --from=builder /bin/irqtop /usr/local/bin/irqtop

ENTRYPOINT ["/usr/local/bin/irqtop"]
CMD ["-interval", "1s"]
