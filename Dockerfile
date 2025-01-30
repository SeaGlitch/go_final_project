FROM golang:1.22.5

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 7540

RUN go build -o /my_go_project

CMD ["/my_go_project"]