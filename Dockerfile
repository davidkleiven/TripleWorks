# Stage 1: The builder stage (use an official Go image)
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache make

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN make build

# Stage 2: The final, minimal scratch image
FROM scratch

COPY --from=builder /app/tripleworks /tripleworks

CMD ["/tripleworks"]

