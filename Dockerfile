FROM golang

RUN mkdir -p /home/app

COPY . /home/app
WORKDIR /home/app
RUN go mod download
RUN go build -o /url-shortener

CMD [ "/url-shortener"]