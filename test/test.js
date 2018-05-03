const Application = require("spectron").Application;
const path = require('path');
const chai = require('chai');
const should = chai.should();
const chaiAsPromised = require('chai-as-promised');
const testPage = require('./test.wallet.js');
const electronPath = require('electron-prebuilt')

var page = new testPage();

// var electronPath = path.join(__dirname, '..', 'node_modules', '.bin', 'electron');



// Path to your application
// var appPath = path.join(__dirname,'..','electron/release/linux-unpacked/skycoin');
var appPath = path.join(__dirname,'..','electron/src/electron-main.js');

var app = new Application({
    path: electronPath,
    args: [appPath]
});

global.before(function () {
    chai.should();
    chai.use(chaiAsPromised);
    page.setApp(app);
});

describe('Test Skycoin', function () {
    beforeEach(function () {
        return app.start();
    });

    afterEach(function () {
        return app.stop();
    });

    it('Opening ', function () {
        // chai.expect("yes").to.equal("no");
    });


    
});