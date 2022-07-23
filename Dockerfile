FROM golang:latest
WORKDIR /go-news
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["make", "run"]