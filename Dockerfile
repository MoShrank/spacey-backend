FROM golang:1.16-alpine

LABEL mantainer="Moritz Eich <hey@moritz.dev"

ENV PORT 8080

RUN mkdir /app
ADD . /app/
WORKDIR /app
CMD ["./main"]

HEALTHCHECK CMD curl --fail http://localhost:8080/ping || exit 1