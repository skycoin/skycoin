var TestWallet = function () {
    var app;

    this.windowCount = 1;
    this.pageTitle = 'Hello World!';
    this.helloText = 'Hello from button!';
    this.buttonId = '#helloButton';
    this.textId = '#helloText';

    this.setApp = function(app) {
        this.app = app;
    }

    this.getWindowCount = function() {
        return this.app.client.waitUntilWindowLoaded().getWindowCount();
    }

    this.getApplicationTitle = function() {
        return this.app.client.waitUntilWindowLoaded().getTitle();
    }

    this.clickButtonAndGetText = function() {
        return this.app.client.click(this.buttonId).getText(this.textId);
    }
}

module.exports = TestWallet;