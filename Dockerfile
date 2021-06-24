FROM pefish/ubuntu-go:v1.16 as builder
WORKDIR /app
ENV GO111MODULE=on
COPY ./ ./
RUN go get -u github.com/pefish/go-build-tool/cmd/...@v0.0.7
RUN make

FROM pefish/ubuntu18_04:v1.2
WORKDIR /app
COPY --from=builder /app/build/bin/linux/ /app/bin/
ENV GO_CONFIG /app/config/pom.yaml
ENV GO_SECRET /app/secret/pom.yaml
CMD ["/app/bin/ether-clef", "--help"]

# docker buildx build --platform linux/amd64 --push -t pefish/ether-clef:v0.0.1 .
