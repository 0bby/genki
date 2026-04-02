FROM node:22 AS frontend

ARG GENKI_VERSION
ENV VITE_GENKI_VERSION=$GENKI_VERSION

WORKDIR /client
COPY ./client/package.json ./client/package-lock.json ./
RUN npm ci
COPY ./client .
RUN npm run build

FROM golang:1.24 AS backend

ARG GENKI_VERSION
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-X main.Version=$GENKI_VERSION" -o app ./cmd/api

FROM debian:bookworm-slim AS final

WORKDIR /app

RUN apt-get update && \
	apt-get install -y ca-certificates && \
	rm -rf /var/lib/apt/lists/*

COPY --from=backend /app/app ./app
COPY --from=frontend /client/build ./client/build
COPY ./client/public ./client/public
COPY ./db ./db

EXPOSE 4110

ENTRYPOINT ["./app"]
