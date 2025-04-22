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
COPY ./cmd ./cmd

RUN go build -o scheduler ./cmd/app/main.go

FROM scratch

LABEL authors="ekvo"

WORKDIR /usr/src/app

COPY ./storage ./storage
COPY ./web ./web
COPY --from=builder /usr/src/build/scheduler /usr/src/app/scheduler

ENV TODO_PORT=8000
ENV TODO_DBFILE=./storage/scheduler.db
ENV TODO_PASSWORD=777524f0cf9c792596eb2b3c57801dbd37b6999910d7e693922ab25c9193faa9
ENV ALGORITHM_TASK_DATE=nextdate
ENV TODO_SECRET_KEY=StatusSeeOther
ENV PATH_DIR_WEB=./web

EXPOSE ${TODO_PORT}

ENTRYPOINT ["/usr/src/app/scheduler"]