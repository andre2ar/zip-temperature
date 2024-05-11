FROM golang:1.22-alpine as build
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/cloudrun .
COPY --from=build /app/app.env.example app.env
ENTRYPOINT ["./cloudrun"]