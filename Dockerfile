FROM golang:1.11 AS builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/learn-opentracing
COPY Gopkg.toml Gopkg.lock ./
COPY . ./

# Download and install the latest release of dep
ADD  https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN dep ensure --vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:3.7

ENV TZ=Asia/Tehran

RUN apk add --update --no-cache tzdata ca-certificates && \
        cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
        echo "${TZ}" > /etc/timezone && rm -rf /var/cache/apk/*

COPY --from=builder /app ./

ENTRYPOINT ["./app"]