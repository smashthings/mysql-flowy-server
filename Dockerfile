FROM golang:1.19 AS builder

COPY . /go/src/scripted.dog/flowy-servers/standalone

WORKDIR /go/src/scripted.dog/flowy-servers/standalone
RUN go mod download && go mod verify
RUN go build -v -o bin/standalone

FROM smasherofallthings/base-image

COPY --from=builder /go/src/scripted.dog/flowy-servers/standalone/bin/standalone /usr/local/standalone
RUN chmod a+x /usr/local/standalone
CMD ["/usr/local/standalone"]

