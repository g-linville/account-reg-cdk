FROM golang:1.20 AS build
COPY / /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build make build

FROM alpine AS release
RUN apk add -U --no-cache ca-certificates
COPY --from=build /src/aws-registration-callback /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/aws-registration-callback"]