ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server/

FROM alpine:3.21 AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    && \
    update-ca-certificates

COPY --from=build /bin/server /bin/

HEALTHCHECK  --interval=30s --start-interval=2s --timeout=3s --start-period=5s --retries=5 \
    CMD wget --no-verbose --spider http://127.0.0.1:8080/health || exit 1

ENTRYPOINT [ "/bin/server" ]
