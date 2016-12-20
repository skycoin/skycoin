import { SecurityContext, SanitizationService } from '../../core_private';
export { SecurityContext };
/** Marker interface for a value that's safe to use in a particular context. */
export interface SafeValue {
}
/** Marker interface for a value that's safe to use as HTML. */
export interface SafeHtml extends SafeValue {
}
/** Marker interface for a value that's safe to use as style (CSS). */
export interface SafeStyle extends SafeValue {
}
/** Marker interface for a value that's safe to use as JavaScript. */
export interface SafeScript extends SafeValue {
}
/** Marker interface for a value that's safe to use as a URL linking to a document. */
export interface SafeUrl extends SafeValue {
}
/** Marker interface for a value that's safe to use as a URL to load executable code from. */
export interface SafeResourceUrl extends SafeValue {
}
/**
 * DomSanitizationService helps preventing Cross Site Scripting Security bugs (XSS) by sanitizing
 * values to be safe to use in the different DOM contexts.
 *
 * For example, when binding a URL in an `<a [href]="someValue">` hyperlink, `someValue` will be
 * sanitized so that an attacker cannot inject e.g. a `javascript:` URL that would execute code on
 * the website.
 *
 * In specific situations, it might be necessary to disable sanitization, for example if the
 * application genuinely needs to produce a `javascript:` style link with a dynamic value in it.
 * Users can bypass security by constructing a value with one of the `bypassSecurityTrust...`
 * methods, and then binding to that value from the template.
 *
 * These situations should be very rare, and extraordinary care must be taken to avoid creating a
 * Cross Site Scripting (XSS) security bug!
 *
 * When using `bypassSecurityTrust...`, make sure to call the method as early as possible and as
 * close as possible to the source of the value, to make it easy to verify no security bug is
 * created by its use.
 *
 * It is not required (and not recommended) to bypass security if the value is safe, e.g. a URL that
 * does not start with a suspicious protocol, or an HTML snippet that does not contain dangerous
 * code. The sanitizer leaves safe values intact.
 */
export declare abstract class DomSanitizationService implements SanitizationService {
    /**
     * Sanitizes a value for use in the given SecurityContext.
     *
     * If value is trusted for the context, this method will unwrap the contained safe value and use
     * it directly. Otherwise, value will be sanitized to be safe in the given context, for example
     * by replacing URLs that have an unsafe protocol part (such as `javascript:`). The implementation
     * is responsible to make sure that the value can definitely be safely used in the given context.
     */
    abstract sanitize(context: SecurityContext, value: any): string;
    /**
     * Bypass security and trust the given value to be safe HTML. Only use this when the bound HTML
     * is unsafe (e.g. contains `<script>` tags) and the code should be executed. The sanitizer will
     * leave safe HTML intact, so in most situations this method should not be used.
     *
     * WARNING: calling this method with untrusted user data will cause severe security bugs!
     */
    abstract bypassSecurityTrustHtml(value: string): SafeHtml;
    /**
     * Bypass security and trust the given value to be safe style value (CSS).
     *
     * WARNING: calling this method with untrusted user data will cause severe security bugs!
     */
    abstract bypassSecurityTrustStyle(value: string): SafeStyle;
    /**
     * Bypass security and trust the given value to be safe JavaScript.
     *
     * WARNING: calling this method with untrusted user data will cause severe security bugs!
     */
    abstract bypassSecurityTrustScript(value: string): SafeScript;
    /**
     * Bypass security and trust the given value to be a safe style URL, i.e. a value that can be used
     * in hyperlinks or `<iframe src>`.
     *
     * WARNING: calling this method with untrusted user data will cause severe security bugs!
     */
    abstract bypassSecurityTrustUrl(value: string): SafeUrl;
    /**
     * Bypass security and trust the given value to be a safe resource URL, i.e. a location that may
     * be used to load executable code from, like `<script src>`.
     *
     * WARNING: calling this method with untrusted user data will cause severe security bugs!
     */
    abstract bypassSecurityTrustResourceUrl(value: string): any;
}
export declare class DomSanitizationServiceImpl extends DomSanitizationService {
    sanitize(ctx: SecurityContext, value: any): string;
    private checkNotSafeValue(value, expectedType);
    private sanitizeHtml(value);
    bypassSecurityTrustHtml(value: string): SafeHtml;
    bypassSecurityTrustStyle(value: string): SafeStyle;
    bypassSecurityTrustScript(value: string): SafeScript;
    bypassSecurityTrustUrl(value: string): SafeUrl;
    bypassSecurityTrustResourceUrl(value: string): SafeResourceUrl;
}
