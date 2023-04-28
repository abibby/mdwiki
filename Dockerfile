FROM golang:1.20 as build

WORKDIR /go/src/github.com/abibby/mdwiki

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build

# Now copy it into our base image.
FROM alpine

RUN apk add --no-cache libc6-compat 
COPY --from=build /go/src/github.com/abibby/mdwiki/mdwiki /mdwiki

VOLUME ["/data"]

CMD ["/mdwiki", "serve", "/data"]
