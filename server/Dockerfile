FROM golang:1.17.7-bullseye AS server-build
RUN apt-get update && apt-get install -y unzip

WORKDIR /build/

ADD https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip protoc.zip
RUN unzip protoc.zip -d protoc
ENV PATH="/build/protoc/bin:${PATH}"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

COPY ./ ./
RUN make

FROM debian:bullseye
COPY --from=server-build /build/lovelove ./
COPY --from=server-build /build/proto/lovelove.proto ./proto/
ENTRYPOINT /lovelove
