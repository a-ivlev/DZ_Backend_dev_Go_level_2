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
	kubectl -n myuser apply -f ./staticsrv/k8s

port forwarding:
	kubectl -n myuser port-forward deployment/statics 8000:8080

delete:
	kubectl -n myuser delete -f ./staticsrv/k8s

busyboxplus:
	kubectl -n myuser run curl --image=radial/busyboxplus:curl -i --tty --rm