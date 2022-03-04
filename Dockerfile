# Build Stage
FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache git 

ENV USER=appuser
ENV UID=10001

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod verify


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o api

# Output Stage
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/api /go/bin/api

USER appuser:appuser

EXPOSE 3000

ENTRYPOINT [ "/go/bin/api" ]
