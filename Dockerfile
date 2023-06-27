FROM golang:1.20-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/olympus

RUN go build

FROM alpine

WORKDIR /app

COPY --from=build /app/cmd/olympus /app/olympus

ENTRYPOINT ["./olympus"]
