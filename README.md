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
