FROM golang:1.20 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o authopia main.go

FROM alpine:3.12
WORKDIR /app
COPY --from=build /app/authopia .
CMD ["./authopia"]
