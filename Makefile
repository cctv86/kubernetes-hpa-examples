COMPONENTS?=kubernetes-hpa-examples
BINARYNAME?=app
TAG?=latest


.PHONY : clean
clean:
	rm -f app

.PHONY : build
build:
	CGO_ENABLED=1 GOOS=linux go build -o ${BINARYNAME} main.go

.PHONY : image
image:
	docker build . -t prodan/${COMPONENTS}:${TAG}

.PHONY : run
run:
	CGO_ENABLED=1 GOOS=linux go run main.go