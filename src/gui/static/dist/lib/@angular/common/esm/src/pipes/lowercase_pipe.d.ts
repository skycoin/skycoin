import { PipeTransform } from '@angular/core';
/**
 * Transforms text to lowercase.
 *
 * ### Example
 *
 * {@example core/pipes/ts/lowerupper_pipe/lowerupper_pipe_example.ts region='LowerUpperPipe'}
 */
export declare class LowerCasePipe implements PipeTransform {
    transform(value: string): string;
}
