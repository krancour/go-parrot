FROM quay.io/deis/lightweight-docker-go:v0.3.0
ARG BASE_PACKAGE_NAME
ARG LDFLAGS
ENV CGO_ENABLED=0
WORKDIR /go/src/$BASE_PACKAGE_NAME/
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/
RUN for f in $(ls cmd); do go build -o bin/$f -ldflags "$LDFLAGS" ./cmd/$f; done

FROM scratch
ARG BASE_PACKAGE_NAME
ENV PATH=/examples:$PATH
COPY --from=0 /go/src/$BASE_PACKAGE_NAME/bin/ /examples/
