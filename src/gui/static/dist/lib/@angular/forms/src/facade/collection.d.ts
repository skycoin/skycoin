/**
 * Wraps Javascript Objects
 */
export declare class StringMapWrapper {
    static merge<V>(m1: {
        [key: string]: V;
    }, m2: {
        [key: string]: V;
    }): {
        [key: string]: V;
    };
    static equals<V>(m1: {
        [key: string]: V;
    }, m2: {
        [key: string]: V;
    }): boolean;
}
/**
 * A boolean-valued function over a value, possibly including context information
 * regarding that value's position in an array.
 */
export interface Predicate<T> {
    (value: T, index?: number, array?: T[]): boolean;
}
export declare class ListWrapper {
    static removeAll<T>(list: T[], items: T[]): void;
    static remove<T>(list: T[], el: T): boolean;
    static equals(a: any[], b: any[]): boolean;
    static flatten<T>(list: Array<T | T[]>): T[];
}
export declare function isListLikeIterable(obj: any): boolean;
export declare function areIterablesEqual(a: any, b: any, comparator: (a: any, b: any) => boolean): boolean;
export declare function iterateListLike(obj: any, fn: (p: any) => any): void;
