ARG PROJECT
ARG APP_NAME
ARG VERSION
ARG GIT_COMMIT

FROM golang:1.17.2 as builder

ARG PROJECT
ENV PROJECT=$PROJECT

ARG APP_NAME
ENV APP_NAME=$APP_NAME

ARG VERSION
ENV VERSION=$VERSION

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT


ENV GOSUMDB=off
ENV GO111MODULE=on
ENV WORKDIR=${GOPATH}/src/${APP_NAME}

COPY ./staticsrv ${WORKDIR}
WORKDIR ${WORKDIR}/cmd/staticsrv

RUN set -xe ;\
go install -ldflags="-X ${PROJECT}/version.Version=${VERSION} -X ${PROJECT}/version.Commit=${GIT_COMMIT}" ;\
ls -lhtr /go/bin/


FROM golang:1.15.1
EXPOSE 8080
WORKDIR /go/bin
COPY --from=builder /go/bin/${APP_NAME} .
#COPY --from=builder ${GOPATH}/src/k8s-go-app/config/*.env ./config/
ENTRYPOINT ["/go/bin/staticsrv"]