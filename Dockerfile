FROM golang

COPY . .

RUN go get github.com/lib/pq

ENTRYPOINT ["/"]

EXPOSE 8080