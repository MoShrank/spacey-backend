FROM golang:1.17-alpine

LABEL maintainer="Moritz Eich <hey@moritz.dev>"

RUN apk --no-cache add curl

RUN mkdir /app

COPY /bin/ /app/

WORKDIR /app

RUN ls | grep "main"

CMD ["./main"]