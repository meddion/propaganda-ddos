build:
	go build -o antirus -v .

run:
	go build -o antirus -v . && ./anti-rusnya-ddos

curl_proxy_example:
	curl -v https:\/\/www.nalog.gov.ru\/ -x 46.3.150.197:8000 -U 0ShxVd:409mML
