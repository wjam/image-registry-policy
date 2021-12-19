FROM golang:1.17 as build

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make

FROM scratch

COPY --from=build /src/bin/policy /

ENTRYPOINT ["/policy"]
