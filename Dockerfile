FROM node:20-alpine AS build-node

COPY --from=golang:1.20-alpine /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY . .

WORKDIR /app/webapp

RUN npm install

RUN npm run ng run olympus:app-shell:production

FROM golang:1.20-alpine as build-golang

WORKDIR /app

COPY . .

RUN go mod download

RUN go build

FROM alpine

WORKDIR /app

COPY --from=build-node /app/webapp/dist/olympus/browser /app/webapp/dist/olympus/browser

COPY --from=build-golang /app/olympus /app/olympus

ENTRYPOINT [ "./olympus" ]
