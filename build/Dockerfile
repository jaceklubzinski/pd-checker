FROM golang:1.20-alpine as builder

RUN apk update
RUN apk add \
  g++ \
  git \
  musl-dev \
  go \
  tesseract-ocr-dev

WORKDIR /app
COPY go.mod ./

RUN apk add build-base
COPY . /app

RUN go build -o pd-checker .

# Build a small image
FROM alpine

COPY --from=builder /app/pd-checker /app/pd-checker

ENTRYPOINT ["/app/pd-checker"]

CMD ["service", "server"]
