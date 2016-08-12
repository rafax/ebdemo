FROM alpine

ADD ebdemo /go/bin/

EXPOSE 3000

CMD ["/go/bin/ebdemo"]
