FROM golang

WORKDIR $GOPATH/src/ulr_shorten

COPY . .

RUN go mod tidy

RUN go build -o main .

EXPOSE 8080

ARG DB_PASS

CMD ["./main"]