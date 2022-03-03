build:
	go build -o antiprop -v .

docker_build:
	docker build -t antiprop .

docker_push:
	docker build -t antiprop .
	docker tag antiprop:latest lovefromukraine/antiprop:latest
	docker push lovefromukraine/antiprop:latest

curl_proxy_example:
	curl -v https://gebank.ru/ -x 5.157.131.149:8409 -U spiznxfg:r6daod3mfgkz
