FROM golang:1.23.0 AS builder

LABEL stage=builder

ENV CGO_ENABLED=0

ENV GOOS=linux

ENV GOARCH=amd64

WORKDIR /usr/src/build

ADD go.mod ./

ADD go.sum ./

RUN go mod download

COPY ./internal ./internal

COPY ./pkg ./pkg

COPY main.go ./

RUN go build -o scheduler main.go

FROM scratch

LABEL authors="ekvo"

WORKDIR /usr/src/app

COPY ./init ./init

COPY ./storage ./storage

COPY ./web ./web

COPY --from=builder /usr/src/build/scheduler /usr/src/app/scheduler

EXPOSE 8000

ENTRYPOINT ["/usr/src/app/scheduler"]