Angular QR Code
===============

    <qrcode></qrcode>

An AngularJS directive to creates QR Codes using Kazuhiko Arase’s [qrcode-generator](https://github.com/kazuhikoarase/qrcode-generator) library.

[See it in action](http://monospaced.github.io/angular-qrcode).

Usage
-----

as element

    <qrcode data="string"></qrcode>

with options

    <qrcode version="2" error-correction-level="M" size="200" data="string"></qrcode>

with expression, observe changes

    <qrcode data="{{var}}"></qrcode>

Options
-------

Permitted values

* version: 1-10

* error-correction-level: 'L', 'M', 'Q', 'H'

* size: integer

Install
-------

    bower install monospaced/angular-qrcode

Include the [qrcode generator library](https://raw.github.com/monospaced/bower-qrcode-generator/master/js/qrcode.js) and the `qrcode.js` script provided by this component in your app, and add `monospaced.qrcode` to your app’s dependencies.

Demo
----------------

[monospaced.github.io/angular-qrcode](http://monospaced.github.io/angular-qrcode)

Reference
----------------

[QR Code versions](http://www.qrcode.com/en/about/version.html)

[QR Code error correction](http://www.qrcode.com/en/about/error_correction.html)
