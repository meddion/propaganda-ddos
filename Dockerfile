FROM golang:1.17-alpine
 
RUN mkdir /anti-rusnya-ddos
 
COPY . /anti-rusnya-ddos
 
WORKDIR /anti-rusnya-ddos
 
RUN go build -o antirus . 
 
ENTRYPOINT ["/anti-rusnya-ddos/antirus"]