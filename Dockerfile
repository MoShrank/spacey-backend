FROM golang:1.17-alpine

LABEL maintainer="Moritz Eich <hey@moritz.dev"

ENV PORT 8080

RUN mkdir /app

COPY /bin/ /app/

WORKDIR /app

RUN ls | grep "main"

CMD ["./main"]

HEALTHCHECK CMD curl --fail http://localhost:8080/ping || exit 1