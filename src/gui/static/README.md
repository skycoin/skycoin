### Setup Instructions

*NOTE:* This starter kit assumes that you already have bower (http://bower.io/) installed locally. If you don't, then run the following command first: ```npm install -g bower```

1) Node Modules and Bower Components are not included in this repository to keep it light weight. After cloning or pulling changes from this repository, make sure to run the following commands in terminal:

```npm install``` and ```bower install``` in that order.

2) Once everything is installed all you have to do is run ```gulp build``` and your new server will be running at ```http://localhost:5000``` (you can edit the port in the gulpFile). To speed up gulp times, the standard ```gulp``` task does not include copying over static files. Using the standard ```gulp``` task will be useful for most cases, but if you need to rebuild the whole ```dist``` folder, use ```gulp build```.


### Working with this application structure
1) All pipeline, automation, and testing dependencies are in the ```node_modules``` folder (installed using npm), while all third party application libraries are located in the ```libs``` folder (installed using bower).

2) Any additional third party modules and plugins should always be installed automatically whenever possible using ```npm install module_name``` or ```bower install module_name``` with the ```--save``` or ```--save-dev``` suffixes to save the dependencies in the ```package.json``` and ```bower.json``` files.

3) All development takes place in the ```app``` folder. Production files are generated with gulp automatically and pushed to the ```dist``` folder (it will automatically be created the first time the ```gulp``` task is run in terminal post-installation).

4) The ```gulpFile.js``` is clearly commented, defining each task that takes place during pipeline automation. Every file change is watched and new files are automatically pushed to the ```dist``` folder. All files are concatenated into individual files for use on production servers.


### Routes, Controllers and TemplateURLs
NOTE: When creating controllers and services/factories, always follow the proper naming convention of starting with an uppercase letter. Everything else can use camelCase.

1) Default AngularJS applications tend to use the ```angular-route``` plugin that makes use of a main ```ng-view``` directive in the index.html file and standard ```href``` tags for links. This application is using the ```angular-ui-router``` plugin for better route nesting and greater customizability. It makes use of a main ```ui-view``` directive instead of ```ng-view``` and uses an ```sref``` tag for links instead of the normal ```href``` tag. Check out the official documentation for more details: https://github.com/angular-ui/ui-router

2) Due to the modularity of this application structure, standard routing parameters aren't being used. In most examples, routes make use of ```TemplateURL``` and ```controller``` like so:

```
$stateProvider
    .state('home', {
      url: '/',
      templateUrl: './modules/home/home.html',
      controller: './modules/home/homeController.js'
    })
...
```
In this application, each module is set up as an injectible directive with its own controller. So instead of the above example, the home module has a directive called ```homeView``` that can be injected into the HTML like this:
```<div home-view></div>``` (camelcased directives always have to be changed to dashed names when in the HTML). As such, our route config makes use of the ```template``` parameter instead of ```templateURL```. So the routes look like this instead:

```
$stateProvider
    .state('home', {
      url: '/',
      template: '<div home-view></div>'
    });
```
As you can see, it's simpler and cleaner, calling only an HTML ```<div></div>``` tag as a template and leaving everything else contained within the module. This way, if anything changes in the file structure, the routes won't need to be updated.

As we add more options and configuration to each state, further changes to the $stateProvider function becomes necessary, so the current configuration looks like this:

```
var home = {
    name: 'home',
    url: '/',
    template: '<div home-view></div>'
};
var module2 = {
    name: 'module2',
    url: '/module2',
    template: '<div module2-view></div>'
};

$stateProvider.state(home);
$stateProvider.state(module2);
```

With this approach, it's very easy to keep every state object clean and easy to understand.

### Adding Modules
1) Create a new folder in the ```app/modules/``` directory with the following files:

```
index.js
moduleName.html
moduleName.less
moduleNameController.js
moduleNameDirective.js
moduleNameConfig.js (this file is only necessary if you'll be adding sub-modules)
```

2) Change the file contents accordingly. Follow the ```app/modules/home``` files as reference. Make sure to change the naming convention in each file.

3) Add a new state to the ```app/config.js``` file like so:

```
var home = {
    name: 'home',
    url: '/',
    template: '<div home-view></div>'
};
var moduleName = {
    name: 'moduleName',
    url: '/moduleName',
    template: '<div moduleName-view></div>'
};

$stateProvider.state(home);
$stateProvider.state(moduleName);
```

4) Open ```app.js``` and add a requirement for the new module. Make sure to require the entire module folder (browserify will look for the index.js file and use that file as the entry point for all other module dependencies).

```
require('./modules/moduleName').name
```

Your end result should look something like this:
```
'use strict';

require('angular');

module.exports = angular.module('myApp',
	[
		require('./common/common').name,
		require('./modules/moduleName').name
	])
	.config(require('./appConfig'));
```

After those steps are complete, you should be able to see the contents of your new module at the URL you specified in step 3.

NOTE: This same process applies to sub-modules, except you will treat the module directory as the root path, create a ```moduleConfig.js``` file where you will define module-specific states and options, and then require the sub-module in the module's ```index.js``` file. You could actually do this with the main ```modules``` directory, and use it to "require" all of your modules instead of app.js and simply call ```require('./modules').name``` instead of ```require('./modules/moduleName').name```. It's all up to you and how deep you want to go with the modularity.

### Adding Third Party Vendor JS and CSS files to the pipeline
Instead of bloating the index.html file with a list of scripts and link tags, all CSS and Javascript files from Vendors are bundled and concatenated into single ```vendor.css``` and ```vendor.js``` files using the Gulp pipeline. To add vendor files to your workflow, all you have to do is access the ```Gulpfile.js``` and add the relative path to the file from the ```libs``` directory to the appropriate location in the *"File Paths"* section.

For CSS files, add the path to the *VendorCSS* workflow.
For JS files, add the path to the *VendorJS* workflow.

NOTE: This is meant strictly for third party libraries that cannot be installed using ```npm install``` or ```bower install```. You should use one of those two methods primarily for installing third party libraries so that you can easily inject them into your modules.

### Learning Resouces
- https://github.com/curran/screencasts/tree/gh-pages/introToAngular
- https://www.codeschool.com/courses/shaping-up-with-angular-js
- http://egghead.io
- http://thinkster.io
