FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o aws-sdk-test .
RUN echo "Successfully sent from kubernetes" > mytest.txt
CMD ["/app/aws-sdk-test"]