# STEP 1 build executable binary

FROM golang:alpine as builder

RUN apk add git

WORKDIR $GOPATH/src/sensit-callback/
ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GO111MODULE=on

COPY /go.mod /go.sum $GOPATH/src/sensit-callback/

RUN go version && \
    go mod download

COPY / $GOPATH/src/sensit-callback/

#build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o $GOBIN/sensit-callback .

# STEP 2 build a small image

# start from scratch
FROM scratch

# Copy our static executable
COPY --from=builder $GOBIN/sensit-callback /sensit-callback
ENTRYPOINT ["/sensit-callback"]
