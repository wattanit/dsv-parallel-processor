##
## BUILDER
##

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /dsv-parallel-processor

##
## DEPLOY
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /dsv-parallel-processor /dsv-parallel-processor

ENTRYPOINT [ "/dsv-parallel-processor" ]