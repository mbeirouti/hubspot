FROM golang:1.13.4

RUN mkdir /usr/local/go/src/Hubspot

COPY . /usr/local/go/src/Hubspot

WORKDIR /usr/local/go/src/Hubspot

RUN go build -o main

RUN chmod +x main

CMD ["./main"]
