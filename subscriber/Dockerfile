FROM scratch

ADD ../ca-certificates.crt /etc/ssl/certs/

# Run this in Terminal first:
# CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o subscriber

ADD subscriber /

ENTRYPOINT ["/subscriber"]
