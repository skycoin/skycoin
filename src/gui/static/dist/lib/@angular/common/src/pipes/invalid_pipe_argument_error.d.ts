/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { Type } from '@angular/core';
import { BaseError } from '../facade/errors';
export declare class InvalidPipeArgumentError extends BaseError {
    constructor(type: Type<any>, value: Object);
}
