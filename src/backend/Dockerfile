FROM golang:1.23.3

WORKDIR /app/backend

COPY . .

RUN go build -o main .

RUN chmod +x ./main

EXPOSE 4000

CMD ["./main"]

