FROM golang

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY ./app .

RUN go build -o cc ./cmd

CMD [ "/cc" ]
