FROM golang:1.25-bookworm AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o issues-2

CMD ["/build/issues-2"]
