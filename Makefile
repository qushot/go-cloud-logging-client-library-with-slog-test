# Evvironment Variables
GOOGLE_CLOUD_PROJECT := your-project-id
REGION := asia-northeast1
REPO_NAME := your-project-id
SERVICE_NAME := your-service-name
SERVICE_IDENTIFIER := your-service-identifier

IMAGE_NAME := $(REGION)-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/$(REPO_NAME)/$(SERVICE_NAME):latest
URL := https://$(SERVICE_NAME)-$(SERVICE_IDENTIFIER)-an.a.run.app/

# Container
.PHONY: container/init container/build container/push container/deploy
container/init:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev
container/build:
	docker build -t $(IMAGE_NAME) .
container/push:
	docker push $(IMAGE_NAME)
container/deploy:
	gcloud run deploy $(SERVICE_NAME) \
	--image $(IMAGE_NAME) \
	--platform managed \
	--region $(REGION) \
	--allow-unauthenticated \
	--set-env-vars=GOOGLE_CLOUD_PROJECT=$(GOOGLE_CLOUD_PROJECT)
container/all: container/build container/push container/deploy

# Go
.PHONY: go/run
go/run:
	go run main.go

# Request
.PHONY: req/get req/post
req/get:
	curl -X GET $(URL)
req/post:
	curl -X POST -H "Content-Type: application/json" -d '{"message": "hello"}' $(URL)
