# Build the project
FROM golang:1.20 as builder

WORKDIR /go/src/gitlab.com/route-kz/auth-api
ADD . .

RUN make build
#RUN make test

# Create production image for application with needed files
FROM golang:1.20.5-alpine3.18

EXPOSE 8000

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/gitlab.com/route-kz/auth-api .

CMD ["./bin/auth-api"]
