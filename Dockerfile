FROM golang:1.15-alpine

RUN mkdir /app 

WORKDIR /app 

COPY . /app/ 

RUN go build -o main . 

CMD ["/app/main"]