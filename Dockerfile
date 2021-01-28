FROM golang:1.14

COPY gateway .

EXPOSE 8080

CMD ["./gateway"]