#Build images
FROM golang:alpine AS builder
RUN apk update && apk upgrade
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates
RUN addgroup -S myapp && adduser -S -u 10000 -g myapp myapp
WORKDIR /app
COPY . .
WORKDIR /app/receiver
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix cgo -ldflags '-extldflags "-static"' -tags timetzdata -o receiver .

# Deploy image
FROM scratch
COPY --from=builder /app/receiver/receiver .
# copy ca certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 # copy users from builder (use from=0 for illustration purposes)
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Americas/New_York
ENV ZONEINFO=/zoneinfo.zip
USER myapp
CMD ["./receiver"]
