build:
	docker build . --no-cache -t coinscan-image

run:
	docker run --name coinscan-container -p 3000:80 -it coinscan-image

stop:
	docker stop coinscan-container

push: #build
	docker tag coinscan-image alexrondon89/coinscan-image
	docker push alexrondon89/coinscan-image

run-from-hub:
	docker run -p 3000:3000 alexrondon89/coinscan-image
