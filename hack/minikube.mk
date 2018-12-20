 

.PHONY: logs
logs: ##@setup Shows logs.
	ktail -n container-image-builder

.PHONY: minikube-await
minikube-await:
	./hack/await-minikube.sh

.PHONY: minikube-start
minikube-start: ##@minikube start minikube
	sudo CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none --kubernetes-version=v1.12.0

.PHONY: minikube-stop
minikube-stop: ##@minikube stop minikube
	sudo minikube stop

.PHONY: minikube-delete
minikube-delete: ##@minikube remove minikube
	sudo minikube delete
