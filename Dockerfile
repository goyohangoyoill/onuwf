FROM golang:1.17-alpine as builder

WORKDIR /workspace

RUN mkdir onuwf
COPY . /workspace/onuwf

WORKDIR /workspace/onuwf/game
RUN go mod download

WORKDIR /workspace/onuwf/util/json
RUN go mod download

WORKDIR /workspace/onuwf/util
RUN go mod download

WORKDIR /workspace/onuwf
RUN go mod download

RUN go build -o onuwf

FROM alpine

WORKDIR /workspace
RUN mkdir onuwf
COPY --from=builder /workspace/onuwf /workspace/onuwf

RUN chmod +x /workspace/onuwf/refresh.sh
RUN chmod +x /workspace/onuwf/onuwf
WORKDIR /workspace/onuwf

ENTRYPOINT ["/bin/sh","/workspace/onuwf/refresh.sh"]
