FROM node:20.2-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci

COPY . .

RUN npm run ng run olympus:app-shell:production

FROM ghcr.io/atuleu/angular-to-http:latest

WORKDIR /app

COPY --from=build /app/dist/olympus/browser /app/webapp/dist/olympus/browser

ENTRYPOINT [ "./angular-to-http", "/app/webapp/dist/olympus/browser" ]
