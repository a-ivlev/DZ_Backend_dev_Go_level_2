
docker build:
	docker build -f ./Dockerfile -t docker.io/aivlev/statics:v1.0.0 .

docker push:
	docker push docker.io/aivlev/statics:v1.0.0

docker run:
	docker run --name statics -p 8080:8080 aivlev/statics:v1.0.0

docker start:
	docker start statics

docker stop:
	docker stop statics

start:
	minikube start

stop:
	minikube stop

create namespace:
	kubectl create namespace myuser

apply:
	kubectl -n myuser apply -f ./k8s

portforwarding:
	kubectl -n myuser port-forward deployment/statics 8000:8080

delete:
	kubectl -n myuser delete -f ./k8s

busyboxplus:
	kubectl -n myuser run curl --image=radial/busyboxplus:curl -i --tty --rm

events:
	kubectl -n myuser get events

dashboard:
	minikube dashboard


#initContainer v1.0.0
applyinitContainer:
	kubectl -n myuser apply -f ./k8s/initContainer/

portforwarding:
	kubectl -n myuser port-forward deployment/statics 8000:8080

deleteinitContainer:
	kubectl -n myuser delete -f ./k8s/initContainer/


#write here path for your project
PROJECT :=

GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := v1.0.1
APP_NAME := staticssrv

all: run

run:
	#cd staticssrv/cmd/staticssrv
	go install -ldflags="-X '$(PROJECT)/version.Version=$(VERSION)' -X '$(PROJECT)/version.Commit=$(GIT_COMMIT)'" && $(APP_NAME)

dockerbuild:
	docker build -f ./Dockerfile -t docker.io/aivlev/$(APP_NAME):$(VERSION) .

dockerpush:
	docker push docker.io/aivlev/$(APP_NAME):$(VERSION)