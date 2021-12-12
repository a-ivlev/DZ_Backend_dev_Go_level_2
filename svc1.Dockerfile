FROM golang:1.17.2 as builder

ENV GOSUMDB=off
ENV GO111MODULE=on
ENV WORKDIR=${GOPATH}/src/

WORKDIR ${WORKDIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./red ${WORKDIR}/red

COPY ./svc1 ${WORKDIR}

# Запускаем компиляцию программы на go и сохраняем полученный бинарный файл server в директорию /rental/ образа.
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../rental/server .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ${GOPATH}/bin/router .


FROM golang:1.15.1

EXPOSE 80

WORKDIR /go/bin
COPY --from=builder /go/bin/router .

ENTRYPOINT ["/go/bin/router"]