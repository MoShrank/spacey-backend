FROM golang:1.17-alpine

LABEL maintainer="Moritz Eich <hey@moritz.dev"

RUN mkdir /app

COPY ./services/user-service/bin/ /app/

WORKDIR /app

RUN ls | grep "main"

EXPOSE 80

CMD ["./main"]