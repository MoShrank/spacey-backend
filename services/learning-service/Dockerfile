FROM golang:1.17-alpine

LABEL maintainer="Moritz Eich <hey@moritz.dev>"

RUN apk --no-cache add curl

RUN mkdir /app

COPY /bin/ /app/

WORKDIR /app

RUN ls | grep "main"

CMD ["./main"]

HEALTHCHECK CMD curl --fail http://localhost:80/ping || exit 1