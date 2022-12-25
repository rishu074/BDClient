
# The BDclient

This is the client for BDagent

## Installation

This agent is only built for ubuntu/linux based distributions.

```bash
# Create the directory for agent
$ mkdir /var/apps
$ mkdir /var/apps/bdclient
$ cd /var/apps/bdclient
$ mkdir logs

# Installing required files
$ sudo apt update -y && sudo apt upgrade -y
$ sudo apt install nginx zip unzip wget curl git -y

# Install the latest binaries
$ wget https://github.com/NotRoyadma/BDClient/releases/latest/download/client
$ wget https://github.com/NotRoyadma/BDClient/releases/latest/download/bdclient.service
```

### Setup configuration file

```bash
# Create a config file
$ touch config.yml
```

### **IMPORTANT STEP**
The data structure is important to know for you, so, what is data structure.
as you can see in the format of config file there is a `DataDirectory` parameter.
This is the data directory in which the data should be structured like this in order for bdsystem
to work.

### Data structure
In this example we will take the `DataDirectory` as `/uploads` and `DataFile` as `data.zip`.

#### File structure
```
 uploads
 ├── data-folder1
 │   └── data.zip
 ├── data-folder2
 │   └── data.zip
 ├── data-folder3
 │   └── data.zip
 └── data-folder4
     └── data.zip
```
## **The file structure should be strictly in this manner, or the system will not gonna work**

#### Data format of config file

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Name`      | `string` | **Required**. Name of the application |
| `Version`      | `string` | **Don't change it**. |
| `port`      | `integer` | **Required**, Port of the application |
| `node`      | `string` | **Required**, The current uploader's name  |
| `DataDirectory`      | `string` | **Required**, The path where the data to send is.  |
| `DataFile`      | `string` | **Required**, The file which is zipped  |
| `Token`      | `string` | **Required**, The authorization token for BDAgent  |
| `ServiceToken`      | `string` | **Required**, The local token of the api to hit  |
| `IpHeader`      | `string` | **Required**, leave it `default` or if using cloudflare change it  |
| `remote`      | `string` | **Required**, The remote url on which bdagent is hosted.  |
| `ssl`      | `boolean` | **Required**, If the remote is ssl or not.  |

An example demostration to config file
```
Name: "BDClient"
Version: "1.0.0"
port: 1337
node: "game1"
DataDirectory: "./uploads"
DataFile: "data.zip"
Token: "SomerandomToken"
ServiceToken: "mysupersecretToken"
IpHeader: "default"
remote: "localhost:1337"
ssl: true
```
### Setting up Systemd service
```bash
# Copy the service to systemd directory
$ cd /var/apps/bdclient
$ mv bdclient.service /etc/systemd/system/
$ systemctl enable --now bdclient.service
```

### There are prebuilt loggers for http and application
```bash
# To view the app logs
$ cat /var/apps/bdclient/logs/app.log

# To view the http logs
$ cat /var/apps/bdclient/logs/http.log

# To view the error logs (if any)
$ cat /var/apps/bdclient/logs/app.error.log

# You can also view the systemctl service status by doing
$ systemctl status bdclient.service

# You can view live http,app and error logs by doing
$ journalctl -u bdclient.service -e --follow
```


## Authors

- [@NotRoyadma](https://www.github.com/NotRoyadma)

