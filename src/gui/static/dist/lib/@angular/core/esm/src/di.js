/**
* @module
* @description
* The `di` module provides dependency injection container services.
*/
export { InjectMetadata, OptionalMetadata, InjectableMetadata, SelfMetadata, HostMetadata, SkipSelfMetadata, DependencyMetadata } from './di/metadata';
// we have to reexport * because Dart and TS export two different sets of types
export * from './di/decorators';
export { forwardRef, resolveForwardRef } from './di/forward_ref';
export { Injector } from './di/injector';
export { ReflectiveInjector } from './di/reflective_injector';
export { Binding, ProviderBuilder, bind, Provider, provide } from './di/provider';
export { ResolvedReflectiveFactory, ReflectiveDependency } from './di/reflective_provider';
export { ReflectiveKey } from './di/reflective_key';
export { NoProviderError, AbstractProviderError, CyclicDependencyError, InstantiationError, InvalidProviderError, NoAnnotationError, OutOfBoundsError } from './di/reflective_exceptions';
export { OpaqueToken } from './di/opaque_token';
//# sourceMappingURL=di.js.map