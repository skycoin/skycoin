# Skycoin desktop client

The Skycoin wallet ships with a web interface which can be ran from the browser and/or Electron.

The project contains both the source (src) and target (dist) files of this web interface.

## Prerequisites

The Skycoin web interface requires Node 8.10.0 or higher, together with NPM 5.6 or higher.

## Installation

This project is generated using Angular CLI, therefore it is adviced to first run `npm install -g @angular/cli`.

Dependencies are managed with NPM 5, to install these run `npm install`.

You will only have to run this again, if any dependencies have been changed in the `package-lock.json` file.

## Compiling new target files

After pulling the latest code, you might first have to update your dependencies, in case someone else has updated them. 
You should always do this when compiling new production files:

```
rm -rf node_modules
npm install 
```

This will remove the current dependencies, and install them from the `package-lock.json`.

To compile new target files, you will have to run: `npm run build`.

## Development server

Run `npm start` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

Please note that you will most likely receive CORS errors as there's a difference between the port number of the source and destination.

As a work-around, the development server will create a proxy from `http://localhost:4200/api` to `http://127.0.0.1:6420/`.

You can route all calls to this address by changing the url property on the ApiService class.

## Swaplab integration

The wallet has a section for buying coins using the Swaplab service and, during development, it is possible to configure its operation,
in order to test its different parts without having to make a real exchange. The configuration is done by modifying the properties of
the `swaplab` object in [environment.ts](/src/gui/static/src/environment/environment.ts).

If you set `swaplab.activateTestMode` to `true`, all operations made with the Swaplab integration will be faked. When creating a new
operation or opening an old one using the history, the system will show the operation progressing every few seconds, from the state in
which the user is expected to make the initial deposit to the moment in which the operation is completed.
If `swaplabTests.activateTestMode` is `false`, the integration will work as it should and you will be able to make real exchange transactions.

Also `swaplab` has a property called `endStatusInError`. If you set it to `true` and `swaplabTests.activateTestMode` is also `true`,
all the faked operations will end in an error, instead of success.

## Purchase API (teller)

Please note that at the moment the Purchase API (teller) is both offline and not supporting CORS headers.

While event.skycoin.com is not working, we will have to run the purchase API locally.

Similar as the solution for the above CORS issue, you can circumvent CORS issues by changing the url property to '/teller/'

## Style guide

As an Angular application, we try to follow the [Angular style guide](https://angular.io/guide/styleguide).

## Translations

You can find information about how to work with translation files in the [Translations README](/src/gui/static/src/assets/i18n/README.md).
