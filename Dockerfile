FROM golang:alpine AS builder

WORKDIR usr/src/authService

#dependencies
COPY go.mod go.sum ./
RUN go mod download

#build
COPY . .
RUN go build -o /usr/local/bin/authService cmd/main.go

FROM alpine AS runner
COPY --from=builder /usr/local/bin/authService /
COPY .env /.env

CMD ["/authService"]