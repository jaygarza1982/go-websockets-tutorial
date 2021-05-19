# # We'll choose the incredibly lightweight
# # Go alpine image to work with
# FROM golang:1.11.1 AS builder

FROM golang:1.16

WORKDIR /app

COPY . .

CMD [ "go", "run", "/app/main.go" ]

# # We want to build and run
# # our application's binary executable
# RUN mkdir /app
# ADD . /app
# WORKDIR /app
# RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

# # the lightweight scratch image we'll
# # run our application within
# FROM alpine:latest AS Production
# COPY --from=builder /app .
# CMD ["./main"]