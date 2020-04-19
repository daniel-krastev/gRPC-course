FROM golang

ENV GOBIN=/go/bin

RUN mkdir -p /go/src/grpc-course/calc && \
    go get github.com/sirupsen/logrus && \
    go get google.golang.org/grpc

WORKDIR /go/src/grpc-course/calc

COPY calc_proto calc_proto
COPY calc_serv calc_serv

RUN go install calc_serv/server.go

EXPOSE 50051

WORKDIR ${GOBIN}

CMD [ "./server" ]