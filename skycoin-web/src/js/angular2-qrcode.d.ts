import { ElementRef, OnChanges, SimpleChanges } from '@angular/core';
export declare class QRCodeComponent implements OnChanges {
    private elementRef;
    ngOnChanges(changes: SimpleChanges): void;
    background: string;
    backgroundAlpha: number;
    foreground: string;
    foregroundAlpha: number;
    level: string;
    mime: string;
    padding: number;
    size: number;
    value: string;
    canvas: boolean;
    constructor(elementRef: ElementRef);
    generate(): void;
}
export declare class QRCodeModule {
}
