
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

## Pterodactyl

If you are using this system for Pterodactyl, then we have a good news for you.
We have already configured the envirnment for Pterodactyl.
>What does that mean?
It means that this project was intended to be built for Pterodactyl software and help Hosting owners
To maintain user's data at one Place, in-case anything wrong happens.

So, if you're using this project for Pterodactyl then follow the instructions.
```bash
# Making bash file to automatic zip data for pterodactyl
$ cd /var/apps/bdclient
$ touch node.sh
```

Now paste the following in `node.sh`.
```bash
#!/bin/bash

# The Node Script 
# - Runs on the pterodactyl's Nodes

# Check if its the root
if [ "$(id -u)" != "0" ]; then
   echo "You should mind running the script as a root user. :)" 1>&2
   exit 1
fi

echo "Starting Backup Script."
# define args 
token=$1
port=$2

# Check if pterodactyl directory exists 
if [ ! -d "/var/lib/pterodactyl/volumes" ]; then
    echo "No pterodactyl/volumes directory."
    exit 1
fi

echo ""
echo "Token $token"
echo "Port $port"
echo ""

if [[ $# -eq 0 ]]; then
    echo "No args specified."
    exit 1
fi

# Check if nginx is install 
if [ -z "$(command -v nginx)" ]; then
    echo "Nginx should be installed to run this script."
    exit 1
elif [ -z "$(command -v zip)" ]; then
    echo "zip should be installed to run this script."
    exit 1
elif [ -z "$(command -v unzip)" ]; then
    echo "unzip should be installed to run this script."
    exit 1
elif [ -z "$(command -v wget)" ]; then
    echo "wget should be installed to run this script."
    exit 1
elif [ -z "$(command -v curl)" ]; then
    echo "curl should be installed to run this script."
    exit 1
fi

# We have nginx + pterodactyl

# Clear uploads folder
mkdir /uploads
cd /uploads
rm -rf *

# Copy all server's data 
echo "Copying data from pterodactyl directory..."
cp -r /var/lib/pterodactyl/volumes/* /uploads/


# Loop through all the folders in /uploads and zip every type of data to data.zip in each folder 
cd /uploads

echo "Looping through servers"
for dir in */
do
        cd $dir
        echo "zipping $dir"
        zip -r data.zip ./*
        shopt -s extglob
        shopt -s dotglob
        for name in *
        do
                if [ "$name" != "data.zip" ] ;
                then
                        echo $name
                        rm -rf $name
                fi
        done
        cd ../;
done

# Now we have exactly what we want in /uploads folder
# Send the request to client
curl -X POST "localhost:$port/upload" -H "token: $token"
```

Now, you are done with the bash file.
>*if you would like to change bash file, you are free to do it*
Now follow the following in-order to setup automated backups

```bash
# Open crontab
$ crontab -e

# paste the following at the last (replace the variables with config values)

0 */12 * * * cd /var/apps/bdclient/ && sudo bash /var/apps/bdclient/node.sh {{ServiceToken}} {{port}}
```

This starts the uploader every 12 hours, creates the Pterodactyl backups and update it to agent.

## Authors

- [@NotRoyadma](https://www.github.com/NotRoyadma)

