FROM golang as base

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get -u "github.com/lib/pq"
RUN go get -u "github.com/imroc/req"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o main .

FROM alpine
COPY --from=base /app/main /main
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk --no-cache add libc6-compat
ENTRYPOINT ["/main"]