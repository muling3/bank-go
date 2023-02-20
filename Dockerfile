FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o bank main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder ./app/bank .
COPY app.env .

EXPOSE 9090
CMD ["/app/bank"]