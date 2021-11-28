
lint:
	golangci-lint run -c golangci-lint.yml

test:
	go test ./...

start:
	minikube start

stop:
	minikube stop

applysrv:
	kubectl apply -f ./kube/srv

applycli:
	kubectl apply -f ./kube/cli

deletesrv:
	kubectl delete -f ./kube/srv

deletecli:
	kubectl delete -f ./kube/cli

busyboxplus:
	kubectl run curl --image=radial/busyboxplus:curl -i --tty --rm
