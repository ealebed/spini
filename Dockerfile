#
# Builder
#
FROM golang:alpine3.15 as builder

WORKDIR /go/src/github.com/ealebed/spini

COPY . /go/src/github.com/ealebed/spini

RUN go build -o bin/spini ./

#
# Runtime
#
FROM alpine:3.15

RUN apk update && apk add git

RUN git config --global user.name "Yevhen Lebid" \
    && git config --global user.email "ealebed@gmail.com"

COPY --from=builder /go/src/github.com/ealebed/spini/bin/spini /bin/spini

WORKDIR /spini

ENTRYPOINT [ "/bin/spini" ]
CMD ["-h"]
