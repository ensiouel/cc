FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./app/ ./

RUN ls

RUN go build -o cc ./cmd

CMD [ "./cc" ]
