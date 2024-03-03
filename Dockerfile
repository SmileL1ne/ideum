FROM golang:1.20.1-alpine3.16 as base
LABEL Team="msarvaro nmagau" Project="Rabbit" 

RUN apk add build-base 
WORKDIR /app 
COPY . .
RUN go build -o forum ./cmd/app/

FROM alpine:3.16
WORKDIR /app
COPY --from=base /app/ /app/

CMD ["./forum"]