"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require("@angular/core");
var QRious = require("qrious");
var QRCodeComponent = (function () {
    function QRCodeComponent(elementRef) {
        this.elementRef = elementRef;
        this.background = 'white';
        this.backgroundAlpha = 1.0;
        this.foreground = 'black';
        this.foregroundAlpha = 1.0;
        this.level = 'L';
        this.mime = 'image/png';
        this.padding = null;
        this.size = 100;
        this.value = '';
        this.canvas = false;
    }
    QRCodeComponent.prototype.ngOnChanges = function (changes) {
        if ('background' in changes ||
            'backgroundAlpha' in changes ||
            'foreground' in changes ||
            'foregroundAlpha' in changes ||
            'level' in changes ||
            'mime' in changes ||
            'padding' in changes ||
            'size' in changes ||
            'value' in changes ||
            'canvas' in changes) {
            this.generate();
        }
    };
    QRCodeComponent.prototype.generate = function () {
        try {
            var el = this.elementRef.nativeElement;
            el.innerHTML = '';
            var qr = new QRious({
                background: this.background,
                backgroundAlpha: this.backgroundAlpha,
                foreground: this.foreground,
                foregroundAlpha: this.foregroundAlpha,
                level: this.level,
                mime: this.mime,
                padding: this.padding,
                size: this.size,
                value: this.value
            });
            if (this.canvas) {
                el.appendChild(qr.canvas);
            }
            else {
                el.appendChild(qr.image);
            }
        }
        catch (e) {
            console.error("Could not generate QR Code: " + e.message);
        }
    };
    return QRCodeComponent;
}());
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "background", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "backgroundAlpha", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "foreground", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "foregroundAlpha", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "level", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "mime", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "padding", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "size", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "value", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Boolean)
], QRCodeComponent.prototype, "canvas", void 0);
QRCodeComponent = __decorate([
    core_1.Component({
        moduleId: 'module.id',
        selector: 'qr-code',
        template: ""
    }),
    __metadata("design:paramtypes", [core_1.ElementRef])
], QRCodeComponent);
exports.QRCodeComponent = QRCodeComponent;
var QRCodeModule = (function () {
    function QRCodeModule() {
    }
    return QRCodeModule;
}());
QRCodeModule = __decorate([
    core_1.NgModule({
        exports: [QRCodeComponent],
        declarations: [QRCodeComponent],
        entryComponents: [QRCodeComponent]
    }),
    __metadata("design:paramtypes", [])
], QRCodeModule);
exports.QRCodeModule = QRCodeModule;
//# sourceMappingURL=angular2-qrcode.js.map