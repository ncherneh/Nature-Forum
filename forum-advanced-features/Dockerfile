FROM golang:latest
COPY . /app
WORKDIR /app
RUN go mod download
RUN apt-get update && apt-get install -y sqlite3
RUN go build -o main 
# EXPOSE 8085
CMD ["/app/main"]