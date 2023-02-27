FROM golang

WORKDIR /build

ADD . ./

RUN go mod download

RUN go build -o /cc ./app/cmd

CMD [ "/cc" ]
