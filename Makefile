build:
	go build -o antirus -v .

curl_proxy_example:
	curl -v https:\/\/www.nalog.gov.ru\/ -x 46.3.150.197:8000 -U 0ShxVd:409mML
