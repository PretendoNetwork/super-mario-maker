# syntax=docker/dockerfile:1
FROM golang:bullseye
EXPOSE 60003/udp
WORKDIR /app
COPY . .
CMD ["go", "run", "."]