/**
* @module
* @description
* Starting point to import all public core APIs.
*/
export * from './src/metadata';
export * from './src/util';
export * from './src/di';
export { createPlatform, assertPlatform, disposePlatform, getPlatform, coreBootstrap, coreLoadAndBootstrap, createNgZone, PlatformRef, ApplicationRef } from './src/application_ref';
export { APP_ID, APP_INITIALIZER, PACKAGE_ROOT_URL, PLATFORM_INITIALIZER } from './src/application_tokens';
export * from './src/zone';
export * from './src/render';
export * from './src/linker';
export { DebugElement, DebugNode, asNativeElements, getDebugNode } from './src/debug/debug_node';
export * from './src/testability/testability';
export * from './src/change_detection';
export * from './src/platform_directives_and_pipes';
export * from './src/platform_common_providers';
export * from './src/application_common_providers';
export * from './src/reflection/reflection';
export { wtfCreateScope, wtfLeave, wtfStartTimeRange, wtfEndTimeRange } from './src/profile/profile';
export { Type, enableProdMode } from "./src/facade/lang";
export { EventEmitter } from "./src/facade/async";
export { ExceptionHandler, WrappedException, BaseException } from "./src/facade/exceptions";
export * from './private_export';
//# sourceMappingURL=index.js.map