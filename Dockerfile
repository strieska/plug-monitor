FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o plug-monitor ./cmd/server


FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/plug-monitor /plug-monitor

EXPOSE 8080

ENTRYPOINT ["/plug-monitor"]