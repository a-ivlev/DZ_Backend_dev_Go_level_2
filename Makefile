
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



#write here your username
USERNAME := aivlev
APP_NAME := staticsrv
VERSION := 1.0.2

#write here path for your project
PROJECT := internal/app
GIT_COMMIT := $(shell git rev-parse HEAD)

build_container:
	docker build --build-arg=PROJECT=$(PROJECT) --build-arg=APP_NAME=$(APP_NAME) --build-arg=VERSION=$(VERSION) --build-arg=GIT_COMMIT=$(GIT_COMMIT) -f ./flag.Dockerfile -t docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) .


push_container:
	docker push docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)