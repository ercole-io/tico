FROM fnproject/go:1.15 as build-stage
WORKDIR /function
WORKDIR /go/src/func/
ENV GO111MODULE=on
COPY . .
COPY config.toml .
RUN cd /go/src/func/ && go build -o tico/
ENTRYPOINT ["./tico"]