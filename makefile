build:
	docker build . --no-cache -t coinscan-image

run:
	docker run --name coinscan-container -p 3000:3000 -it coinscan-image

stop:
	docker stop coinscan-container