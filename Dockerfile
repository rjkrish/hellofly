FROM golang:1.13 AS builder

WORKDIR /go/src/hellofly
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o hellofly

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/hellofly/hellofly /hellofly
RUN mkdir /templates
COPY --from=builder /go/src/hellofly/templates/*.tmpl /templates/
EXPOSE 8080

CMD ["/hellofly"]


