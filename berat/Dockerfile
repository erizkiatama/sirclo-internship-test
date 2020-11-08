# Builder
FROM golang:alpine as builder

LABEL maintainer="Ezra Rizkiatama <erizkiatama@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY . .

RUN go build -o main .

# Distribution
FROM alpine:latest

COPY --from=builder /build/main .
COPY --from=builder /build/.env .
COPY --from=builder /build/views/ ./views/

CMD [ "./main" ]