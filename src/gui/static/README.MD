

If folder dev exists, application will serve web assets out of dev. Otherwise the dist folder will be used.

dev will not be commited to the repo. For development.

dist is the distribution state files

#### 1. Prerequisites

*nodejs* must be installed on your system and the below global node packages must be installed:

- gulp

> npm i -g gulp

- gulp-cli

> npm i -g gulp-cli

- typings

> npm i -g typings@1.3.3

- typescript

> npm i -g typescript@2.0.2

- ts-node

> npm i -g ts-node@1.3.0

Updating NPM
-----

```
sudo apt-get install npm
sudo apt-get install nodejs-legacy
sudo npm cache clean -f

sudo npm install -g n
sudo n stable

node -v
npm -v
```

Setup
-----

```
npm install
npm install --only=dev
```

A folder `node_modules/` should now exist.

Building
----

```
gulp build
```

Will build the assets into the ./dev folder

Updating
----

Copy dev over dist, and commit. To update the repo

