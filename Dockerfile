FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev sqlite

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o jobscraper main.go

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache sqlite

COPY --from=builder /app/jobscraper .

RUN mkdir -p /data

ENV DB_PATH=/data/jobs.db
ENV PORT=8080

EXPOSE 8080

CMD ["./jobscraper"]