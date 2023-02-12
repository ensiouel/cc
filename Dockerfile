FROM golang

ADD . /build

WORKDIR /build/app/cmd

RUN go mod tidy
RUN go mod download
RUN go build -o cc

RUN dir -s

COPY сс /app

WORKDIR /app

ENTRYPOINT ["./сс"]
