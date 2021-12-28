FROM golang:1.17-alpine

LABEL miantainer="Moritz Eich <hey@moritz.dev"

ENV PORT 8080

RUN mkdir /app
ADD . /app/
WORKDIR /app
CMD ["./main"]

HEALTHCHECK CMD curl --fail http://localhost:8080/ping || exit 1