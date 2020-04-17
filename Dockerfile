FROM golang:alpine as builder
RUN apk add git
WORKDIR /gitgopher
ADD . .
RUN go build -ldflags="-s -w" -o /usr/bin/gg .
RUN mkdir -p /.cache/go-build work \
    && chown -R nobody /.cache/go-build /gitgopher
USER nobody
CMD ["gg", "-v"]
#would target scratch, but runtime dependencies for tools on host makes this less resistant.