FROM golang:1.10
EXPOSE 8080

WORKDIR /go/src/github.com/ninjadotorg/handshake-exchange/

RUN ["apt-get","-y", "update"]
RUN curl https://glide.sh/get | sh

# Build source
ADD . /go/src/github.com/ninjadotorg/handshake-exchange/

RUN mkdir ./build
RUN cd ./build
RUN go build -v -o handshake-exchange /go/src/github.com/ninjadotorg/handshake-exchange/main/main.go

# Set timezone.
# ENV TZ=America/Los_Angeles
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Need to init firebase at the beginning (.env doesn't help)
ENV FIREBASE_CREDENTIALS="./credentials/cred.json"

# Launch web.
RUN cd /go/src/github.com/ninjadotorg/handshake-exchange
CMD ["go", "run", "/go/src/github.com/ninjadotorg/handshake-exchange/main/main.go"]
