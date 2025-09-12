# Development

The explorer is divided in 2 parts: the server and the frontend. This file has general development information about both parts.

# Table of Contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Server](#Server)
	- [Server file structure](#server-file-structure)
	- [Formatting](#formatting)
- [Frontend](#frontend)
	- [Frontend file structure](#frontend-file-structure)
  - [Integration with the browser's search functionality](#integration-with-the-browser's-search-functionality)
  - [Page reuse strategy](#page-reuse-strategy)
  - [General configuration and customization](#general-configuration-and-customization)
  - [Translations](#translations)
  - [Dependence with the server](#dependence-with-the-server)
  - [How to start the Angular dev server](#how-to-start-the-angular-dev-server)
  - [Compiling the frontend](#compiling-the-frontend)
  - [Code linting](#code-linting)
  - [e2e Tests](#e2e-Tests)

<!-- /MarkdownTOC -->

## Server

The server is a small program written in Go (Golang) which two functions: serve the contents of the `/dist` folder (the compiled frontend) in [http://127.0.0.1:8001](http://127.0.0.1:8001) and provide an API for getting data from the blockchain. The API provided by the server is just a proxy that gets the data from a Skycoin node.

### Server file structure

The server has only one code file: [explorer.go](explorer.go).

The dependencies have been vendored with `dep`.

### Formatting

`explorer.go` should be formatted with `goimports`. You can do this with:

```sh
make format
```

You must have goimports installed (use `make install-linters`).

## Frontend

The frontend is an Angular application that can be viewed in a web browser. I currently uses Angular 15.

### Frontend file structure

The frontend code is inside the [src](src) folder and it follows the normal structure of an Angular project.

You can find the components, pipes and services (most of the code) in in the [src/app](src/app) folder. Resource files are in [src/assets](src/assets) and the JS libs/tools that were not included with NPM are in [src/js](src/js).

You can find the code of the e2e tests in the [test](test) folder. The unit tests are currently not being used.

### Integration with the browser's search functionality

The explorer is integrated with the browser's search functionality. This means that the user can search for elements of the blockchain directly from the search bar of the browser. The initial configuration is done in [src/environments/environment.ts](src/environments/environment.ts) (development builds) and  [src/environments/environment.ts](src/environments/environment.ts) (production builds). The configuration files are [src/search.dev.xml](src/search.dev.xmll) (development builds) and [src/search.xml](src/search.xml) (production builds).

You can find more information about the format of the configuration files in [opensearch-1-1-draft-6.md](https://github.com/dewitt/opensearch/blob/master/opensearch-1-1-draft-6.md)

### Page reuse strategy

This app has a custom `RouteReuseStrategy` that makes the navigator reuse only the page with information about an address. This means that when changing to a URL (even if only the params were changed), Angular destroys and recreates the page component (custom behavior), except in some cases when navigating to `AddressDetailComponent` (normal Angular behavior).

You can find the custom `RouteReuseStrategy` in [src/app/app.reuse-strategy.ts](src/app/app.reuse-strategy.ts).

### General configuration and customization

The general configuration settings are in [src/app/app.config.ts](src/app/app.config.ts). You can find much more information in [CUSTOMIZATION.md](CUSTOMIZATION.md).

### Translations

You can find information about how to work with translation files in the [Translations README](/src/assets/i18n/README.md).

### Dependence with the server

The frontend needs the Go server for getting the data. The frontend will not normally work well if it is not opened with the Go server or the Angular dev server. This is due to how the connections are made in [src/app/services/api/api.service.ts](src/app/services/api/api.service.ts).

### How to start the Angular dev server

As it's very inconvenient having to build the frontend evey time a change is made, to be able to see it, you can run the frontend using a development server, which shows the changes every time a file is modified. To do so, simply run `npm start` in the root of the repository. After the development server is started, you can access the frontend via [http://localhost:4200](http://localhost:4200).

NOTE: when using the development server you still need to run the Go server and a Skycoin node for the frontend to work.

### Compiling the frontend

To create a production build of the frontend, run:
```sh
make build-ng
```

### Code linting

The code must pass the linting process. To test the code, run:
```sh
npm run lint
```

### e2e Tests

If you are running a Skycoin node normally, the server must be run locally with `npm start` and you can run the e2e tests with:

```sh
npm run e2e
```

If you are running a Skycoin node using the test database (`blockchain-180.db`), the server must be running on `http://127.0.0.1:8001` and you can run the e2e tests with:

```sh
npm run e2e-blockchain-180
```

The second method is the one used in Travis.

Note 1: before running the e2e tests with the test database (`blockchain-180.db`), you must compile the frontend and run the server. The e2e testing process will use the compiled frontend, so if you forget to run the compilation process your lastest changes will be ignored.

Note 2: if the e2e tests fail with an error related to the version of Google Chrome and Chromedriver, you must update the version of Chromedriver on [package.json](package.json) to make it coincide with the version of Google Chrome. After that, run `npm install` to update the dependencies and run the tests again.
