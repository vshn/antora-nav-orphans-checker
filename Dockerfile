# Step 1: build binary
FROM docker.io/library/alpine:3.18 AS build
RUN apk update && apk upgrade && apk add --no-cache go
WORKDIR /app
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY *.go ./
COPY fixture ./fixture
RUN go test
RUN GOOS=linux GOARCH=$TARGETARCH go build

# Step 2: deployment image
FROM docker.io/library/alpine:3.18
WORKDIR /app
COPY --from=build /app/antora-nav-orphans-checker /app/antora-nav-orphans-checker
USER 1001
ENTRYPOINT ["/app/antora-nav-orphans-checker"]
