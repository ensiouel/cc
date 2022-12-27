FROM golang

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY /app .

RUN go build -o app cmd/main.go

CMD [ "/build/app" ]