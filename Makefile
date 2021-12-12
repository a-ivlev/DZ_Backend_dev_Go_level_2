docker_build_router:
	docker build -f ./svc1.Dockerfile -t docker.io/aivlev/router:v1.0.3 .

docker_push_router:
	docker push docker.io/aivlev/router:v1.0.3

docker_run_router:
	docker run --name router -p 8002:80 aivlev/router:v1.0.3

docker_start_router:
	docker start router

docker_stop_router:
	docker stop router




docker_build_acl:
	docker build -f ./svc2.Dockerfile -t docker.io/aivlev/acl:v1.0.1 .

docker_push_acl:
	docker push docker.io/aivlev/acl:v1.0.1

docker_run_acl:
	docker run --name acl -p 8002:80 aivlev/acl:v1.0.1

docker_start_acl:
	docker start acl

docker_stop_acl:
	docker stop acl

