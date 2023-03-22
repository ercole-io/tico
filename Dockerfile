FROM fnproject/go:1.17-dev as build-stage
WORKDIR /function
WORKDIR /go/src/func/
ENV GO111MODULE=on
COPY . .
CMD [ "ls" ]
COPY config.toml .
RUN cd /go/src/func/ && go build -o tico
FROM fnproject/go:1.17-dev
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
ENTRYPOINT ["./tico"]