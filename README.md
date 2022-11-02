<h3 align="center">Bd Client</h3>

---

<p align="center"> Generates Backups and send it to BD Agent
    <br> 
</p>

### Installing

A step by step series of examples that tell you how to get app env running.

First Clone this git repo
Lets assume your installation directory as (DIR)

```
cd (DIR) && mkdir logs
```

Then Configure the config.yml and bdclient.service and there you go

```
mv (DIR)/bdclient.service /etc/systemd/system/ 
systemctl enable --now bdclient.service
```

Then Just Open your crontab
```
crontab -e
```

After all Paste this line there, but make sure your ServiceToken is replaced at (token) and your app port at (aport)

```
0 */12 * * * cd (DIR) && bash (DIR)/node.sh (token) (aport)
```

### Quick installation

If you don't care about which dir to use or else follow theese steps

```
sudo apt update -y && sudo apt upgrade -y
```

Installation of required packages
```
sudo apt install nginx zip unzip wget curl -y
```

Installation of Client
```
mkdir /var/apps && cd /var/apps && mkdir bdclient && cd bdclient && git clone https://github.com/NotRoyadma/BDClient.git .
```

Installing golang
```
cd /usr/local && wget https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
cd /usr/local && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz && rm go1.19.3.linux-amd64.tar.gz
```

Configure
```
mkdir /var/apps/bdclient/logs
```

```
nano /var/apps/bdclient/config.yml
```

Deploy
```
mv /var/apps/bdclient/quick.service /var/apps/ && mv /var/apps/quick.service /etc/systemd/system/bdclient.service
```
```
systemctl enable --now bdclient.service
```
