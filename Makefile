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