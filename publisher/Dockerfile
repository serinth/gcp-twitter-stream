FROM scratch

ADD ../ca-certificates.crt /etc/ssl/certs/

# Run this in Terminal first:
# CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o publisher

ADD publisher /

# Twitter credentials. For workshop purposes only. 
ENV TWITTER_CONSUMER_KEY value
ENV TWITTER_CONSUMER_SECRET value
ENV TWITTER_ACCESS_TOKEN value
ENV TWITTER_ACCESS_SECRET value

ENTRYPOINT ["/publisher", ""]
