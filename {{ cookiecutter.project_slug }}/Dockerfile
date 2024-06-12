# stage: build ---------------------------------------------------------

FROM golang:1.22-alpine as build

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /go/src/github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o bin/{{ cookiecutter.project_slug }} -ldflags "-s -w" github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/cmd

# stage: run -----------------------------------------------------------

FROM alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /go/src/github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/bin/{{ cookiecutter.project_slug }} ./{{ cookiecutter.project_slug }}

ENTRYPOINT ["/app/{{ cookiecutter.project_slug }}"]
