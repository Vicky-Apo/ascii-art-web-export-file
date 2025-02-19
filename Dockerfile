# ------ BUILDER STAGE ------
FROM golang:1.23.4 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ascii-art-web-dockerize .

# ------ FINAL STAGE ------
FROM alpine:3.18
WORKDIR /app

COPY --from=build /app/ascii-art-web-dockerize .
COPY --from=build /app/templates ./templates
COPY --from=build /app/banners ./banners
COPY --from=build /app/static ./static

# Expose the port that the app runs on
EXPOSE 8080

# Metadata
LABEL version="1.0"
LABEL maintainers="kapostolo, vapostol, ikopylov"
LABEL description="Go server for Ascii Art Generator"

# Run the compiled Go application
CMD ["./ascii-art-web-dockerize"]