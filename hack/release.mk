
.PHONY: build-push
build-push: ##@release Build and push the images.
	docker build -t davedamoon/temp1:latest devtools
	docker build -t davedamoon/temp2:latest worker
	docker build -t davedamoon/temp3:latest dispatcher
	docker push davedamoon/temp1:latest
	docker push davedamoon/temp2:latest
	docker push davedamoon/temp3:latest
