FROM golang:1.19.6-alpine3.17 AS install
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS build
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/proxy cmd/proxy/main.go

FROM ubuntu AS setup_ca
WORKDIR /
COPY scripts/gen_cert.sh scripts/v3.ext ./
COPY ca.crt ca.key ./
RUN apt update -y && \
    apt install openssl ca-certificates && \
    mkdir -p /usr/local/share/ca-certificates && \
    cp ca.crt /usr/local/share/ca-certificates/ca.crt && \
    update-ca-certificates

FROM setup_ca AS proxy
WORKDIR /
COPY --from=build /bin/proxy /bin/proxy
CMD ["/bin/proxy"]
