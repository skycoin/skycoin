# Skycoin desktop client

The Skycoin wallet ships with a web interface which can be ran from the browser and/or Electron.

The project contains both the source (src) and target (dist) files of this web interface.

## Prerequisites

The Skycoin web interface requires Node 6.9.0 or higher, together with NPM 3 or higher.

## Installation

This project is generated using Angular CLI, therefore it is adviced to first run `npm install -g @angular/cli`.

To install all Angular, Angular CLI and all other libraries, you will have to run `npm install`.

You will only have to run this again, if any dependencies have been changed in the `package.json` file.

## Compiling new target files

To compile new target files, you will have to run: `npm run build`

## Development server

Run `npm start` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

Please note that you will most likely receive CORS errors as there's a difference between the port number of the source and destination.

As a work-around, the development server will create a proxy from `http://localhost:4200/api` to `http://127.0.0.1:6420/`.

You can route all calls to this address by changing the url property on the ApiService class.

## Purchase API (teller)

Please note that at the moment the Purchase API (teller) is both offline and not supporting CORS headers.

While event.skycoin.net is not working, we will have to run the purchase API locally.

Similar as the solution for the above CORS issue, you can circumvent CORS issues by changing the url property to '/teller/'
