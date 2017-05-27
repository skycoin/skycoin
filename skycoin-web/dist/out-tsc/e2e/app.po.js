"use strict";
var protractor_1 = require('protractor');
var MeanAppPage = (function () {
    function MeanAppPage() {
    }
    MeanAppPage.prototype.navigateTo = function () {
        return protractor_1.browser.get('/');
    };
    MeanAppPage.prototype.getParagraphText = function () {
        return protractor_1.element(protractor_1.by.css('app-root h1')).getText();
    };
    return MeanAppPage;
}());
exports.MeanAppPage = MeanAppPage;
//# sourceMappingURL=app.po.js.map