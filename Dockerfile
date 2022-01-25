FROM golang:1.17.6-alpine3.14

# ENV APP_NAME blobber
# ENV PORT 8080
ENV secret ciao
#${SECRET}
#${APP_NAME}

EXPOSE 8080
WORKDIR /go/src/blobber
COPY go.mod go.sum /go/src/blobber/
RUN go mod download

COPY ./images/ /go/src/blobber/images/ 
COPY ./pages/ /go/src/blobber/pages/
COPY *.go /go/src/blobber/
RUN go build -o blobber

CMD ./blobber
