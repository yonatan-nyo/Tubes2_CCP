FROM golang:1.24.2

WORKDIR /app/scraper

COPY . .

RUN go build -o main .

RUN chmod +x ./main

CMD ["./main"]
