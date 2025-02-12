FROM --platform=$BUILDPLATFORM golang:1.24 AS build

ARG VERSION
ARG REVISION
ARG BUILDTIME
ARG BUILDER

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /src
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags "-X github.com/mgechev/revive/cli.version=${VERSION} -X github.com/mgechev/revive/cli.commit=${REVISION} -X github.com/mgechev/revive/cli.date=${BUILDTIME} -X github.com/mgechev/revive/cli.builtBy=${BUILDER}"

FROM scratch

COPY --from=build /src/revive /revive

ENTRYPOINT ["/revive"]