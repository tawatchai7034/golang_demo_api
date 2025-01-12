build:
	go build \ 
		-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \ 
		-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \ 
		-O app 

image: 
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p:8081:8081 --env-file ./.env -- --link some-mariadb:db \
	--name myapp todo:test