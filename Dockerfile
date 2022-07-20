FROM golang:1.18-alpine AS build
ADD . /goex21
ENV CGO_ENABLED=0
WORKDIR /goex21
RUN go build -o goex21.bin ./cmd/goex21

FROM alpine:latest
COPY --from=build /goex21/goex21.bin /goex21/goex21.bin
EXPOSE 8888
ENTRYPOINT ["/goex21/goex21.bin"]
