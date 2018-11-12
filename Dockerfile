# STEP 1 build executable binary

FROM golang:alpine as builder
WORKDIR $GOPATH/src/sensit-callback/
COPY sensit-callback.go .

RUN apk add git
#get dependancies
#you can also use dep
RUN go get -d -v

#build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o $GOBIN/sensit-callback .

# STEP 2 build a small image

# start from scratch
FROM scratch

# Copy our static executable
COPY --from=builder $GOBIN/sensit-callback /sensit-callback
ENTRYPOINT ["/sensit-callback"]
