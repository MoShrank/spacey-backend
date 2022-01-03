FROM golang:1.17-alpine

LABEL maintainer="Moritz Eich <hey@moritz.dev"

RUN mkdir /app

COPY /bin/ /app/

WORKDIR /app

RUN ls | grep "main"

CMD ["./main"]