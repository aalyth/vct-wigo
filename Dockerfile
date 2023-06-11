FROM golang:alpine

WORKDIR /client
ADD client ./

WORKDIR /app
COPY server/* ./
RUN go mod download
RUN go build -o /app/wigo 

EXPOSE 4000 

CMD ["/app/wigo"]
