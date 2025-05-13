# Start by building the application.
FROM golang:1.24 AS build

WORKDIR /build
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN make build
RUN chmod +x /build/bin/server

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /build/bin/server /
CMD ["./server"]
