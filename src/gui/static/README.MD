

If folder dev exists, application will serve web assets out of dev. Otherwise the dist folder will be used.

dev will not be commited to the repo. For development.

dist is the distribution state files

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

