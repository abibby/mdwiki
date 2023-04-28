FROM golang:1.20 as build

WORKDIR /go/src/github.com/abibby/mdwiki

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

# Now copy it into our base image.
FROM alpine

COPY --from=build /go/src/github.com/abibby/mdwiki/mdwiki /mdwiki

VOLUME ["/data"]

CMD ["/mdwiki", "serve", "/data"]
