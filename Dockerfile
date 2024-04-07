FROM golang:1.20.0 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
ADD . ./
RUN go mod download && make build

FROM build-stage AS run-test-stage
RUN make test

FROM alpine as run-release-stage
WORKDIR /app
COPY --from=build-stage /app/bin ./bin
CMD ["./bin/main"]