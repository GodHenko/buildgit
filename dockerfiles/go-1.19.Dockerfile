FROM golang:1.19-alpine

ENV CODECRAFTERS_DEPENDENCY_FILE_PATHS="go.mod,go.sum"

WORKDIR /app

COPY go.mod go.sum ./

RUN /bin/ash -c "go mod graph | awk '{if ($1 !~ \"@\") print $2}' | xargs -r go get"

RUN apk add --no-cache 'git>=2.40'
