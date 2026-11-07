REGISTRY   := europe-west2-docker.pkg.dev
REPOSITORY := football-predictions-508514/fleps/
TAG        := latest

SERVICES := space command message action

.PHONY: all push $(SERVICES)

all: $(SERVICES)

command: SERVICE_ARG=commandHandler
space: SERVICE_ARG=spaceHandler
message: SERVICE_ARG=messageHandler
action: SERVICE_ARG=actionHandler

$(SERVICES):
	podman build -t $(REGISTRY)/$(REPOSITORY)$@:$(TAG) --build-arg SERVICE=$(SERVICE_ARG) .

push: $(SERVICES)
	for service in $(SERVICES); do \
		podman push $(REGISTRY)/$(REPOSITORY)$$service:$(TAG); \
	done

login:
	gcloud auth print-access-token | podman login -u oauth2accesstoken --password-stdin europe-west2-docker.pkg.dev