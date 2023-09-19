FROM golang:1.19.6-alpine3.17 AS install
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS build
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/api cmd/api/main.go

FROM ubuntu AS setup_ca
WORKDIR /
COPY ca.crt ./
RUN apt update -y && \
    apt install -y ca-certificates && \
    mkdir -p /usr/local/share/ca-certificates && \
    cp ca.crt /usr/local/share/ca-certificates/ca.crt && \
    update-ca-certificates

FROM setup_ca AS api
WORKDIR /
COPY --from=build /bin/api /bin/api
CMD ["/bin/api"]
