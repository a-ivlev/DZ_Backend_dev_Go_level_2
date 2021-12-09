FROM golang:1.17.2 as builder

ENV GOSUMDB=off
ENV GO111MODULE=on
ENV WORKDIR=${GOPATH}/src/

WORKDIR ${WORKDIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./red ${WORKDIR}/red

COPY ./svc2 ${WORKDIR}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ${GOPATH}/bin/acl .


FROM golang:1.15.1

EXPOSE 80

WORKDIR /go/bin
COPY --from=builder /go/bin/acl .

ENTRYPOINT ["/go/bin/acl"]














#FROM golang:1.17.2 as svc2-build
#
## Создаём директорию server и переходим в неё.
#WORKDIR /app
#
## Копируем файлы go.mod и go.sum и делаем загрузку, чтобы вовремя пересборки контейнера зависимости
## подтягивались из слоёв.
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#
### Копируем все файлы из директории ./shortener локальной машины в директорию /app/shortener образа.
#COPY ./svc2 ./svc2
#
## Запускаем компиляцию программы на go и сохраняем полученный бинарный файл server в директорию /rental/ образа.
##RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../rental/server .
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../svc2 ./svc2
#
#
#FROM scratch
#
#WORKDIR /app
#
## Копируем бинарник server из образа builder в корневую директорию.
#COPY --from=svc2-build /svc2 /
#
## Копируем сертификаты и таймзоны.
#COPY --from=server-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=server-build /usr/share/zoneinfo /usr/share/zoneinfo
## Устанавливаем в переменную окружения свою таймзону.
#ENV TZ=Europe/Moscow
#
## Информационная команда показывающая на каком порту будет работать приложение.
#EXPOSE 80
#
### Устанавливаем по дефолту переменные окружения, которые можно переопределить при запуске контейнера.
##ENV SRV_PORT=8035
##ENV SHORTENER_STORE=mem
#
## Запускаем приложение.
#CMD ["/svc2"]