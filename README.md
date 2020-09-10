# laradock-manager	

A simple application for managing laradock containers.	
Developed and tested only on Ubuntu 19.04/18.04	

Made with https://wails.app/ (go & vue.js & vuetify)	

It's still in a proof of concept stage, so...

## Usage

In order to use this your current user need to be able to access docker without sudo.

1.  Create the docker group. `$ sudo groupadd docker`
2. Add your user to the docker group. `$ sudo usermod -aG docker $USER`
3. Log out and log back in so that your group membership is re-evaluated. (or `$ newgrp docker `)


## Build deb file

`$ ./build.sh`


![Preview](https://raw.githubusercontent.com/Lyimmi/laradock-manager/master/assets/laradock-manager-0.4.0.gif)
