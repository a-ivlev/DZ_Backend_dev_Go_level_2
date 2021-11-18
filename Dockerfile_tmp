FROM golang:1.17.2 as server-build

# Создаём директорию server и переходим в неё.
WORKDIR /app

# Копируем файлы go.mod и go.sum и делаем загрузку, чтобы вовремя пересборки контейнера зависимости
# подтягивались из слоёв.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем все файлы из директории ./shortener локальной машины в директорию /app/shortener образа.
COPY ./shortener ./shortener

# Запускаем компиляцию программы на go и сохраняем полученный бинарный файл server в директорию /rental/ образа.
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../rental/server .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../server ./shortener/cmd/


FROM scratch

WORKDIR /app

# Копируем бинарник server из образа builder в корневую директорию.
COPY --from=server-build /server /

# Копируем сертификаты и таймзоны.
COPY --from=server-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server-build /usr/share/zoneinfo /usr/share/zoneinfo
# Устанавливаем в переменную окружения свою таймзону.
ENV TZ=Europe/Moscow

# Информационная команда показывающая на каком порту будет работать приложение.
EXPOSE 8035

# Устанавливаем по дефолту переменные окружения, которые можно переопределить при запуске контейнера.
ENV SRV_PORT=8035
ENV SHORTENER_STORE=mem

# Запускаем приложение.
CMD ["/server"]


FROM golang:1.17.2 as client-build

# Создаём директорию server и переходим в неё.
WORKDIR /app

# Копируем файлы go.mod и go.sum и делаем загрузку, чтобы вовремя пересборки контейнера зависимости
# подтягивались из слоёв.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем все файлы из директории ./shortener локальной машины в текущую директорию (shortener) образа.
COPY ./client ./client

# Запускаем компиляцию программы на go и сохраняем полученный бинарный файл server в директорию / образа.
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../rental/server .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../client ./client/cmd/


FROM scratch

WORKDIR /app

# Копируем бинарник client из образа builder в корневую директорию.
COPY --from=client-build /client /

# Копируем сертификаты и таймзоны.
COPY --from=client-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=client-build /usr/share/zoneinfo /usr/share/zoneinfo
# Устанавливаем в переменную окружения свою таймзону.
ENV TZ=Europe/Moscow

# Информационная команда показывающая на каком порту будет работать приложение.
EXPOSE 8080

# Устанавливаем по дефолту переменные окружения, которые можно переопределить при запуске контейнера.
ENV SRV_PORT=8035
ENV CLI_PORT=8080

# Запускаем приложение.
CMD ["/client"]