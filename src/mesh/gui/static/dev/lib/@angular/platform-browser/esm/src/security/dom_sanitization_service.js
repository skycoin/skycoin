import { sanitizeUrl } from './url_sanitizer';
import { sanitizeStyle } from './style_sanitizer';
import { SecurityContext } from '../../core_private';
import { Injectable } from '@angular/core';
export { SecurityContext };
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
export class DomSanitizationService {
}
export class DomSanitizationServiceImpl extends DomSanitizationService {
    sanitize(ctx, value) {
        if (value == null)
            return null;
        switch (ctx) {
            case SecurityContext.NONE:
                return value;
            case SecurityContext.HTML:
                if (value instanceof SafeHtmlImpl)
                    return value.changingThisBreaksApplicationSecurity;
                this.checkNotSafeValue(value, 'HTML');
                return this.sanitizeHtml(String(value));
            case SecurityContext.STYLE:
                if (value instanceof SafeStyleImpl)
                    return value.changingThisBreaksApplicationSecurity;
                this.checkNotSafeValue(value, 'Style');
                return sanitizeStyle(value);
            case SecurityContext.SCRIPT:
                if (value instanceof SafeScriptImpl)
                    return value.changingThisBreaksApplicationSecurity;
                this.checkNotSafeValue(value, 'Script');
                throw new Error('unsafe value used in a script context');
            case SecurityContext.URL:
                if (value instanceof SafeUrlImpl)
                    return value.changingThisBreaksApplicationSecurity;
                this.checkNotSafeValue(value, 'URL');
                return sanitizeUrl(String(value));
            case SecurityContext.RESOURCE_URL:
                if (value instanceof SafeResourceUrlImpl) {
                    return value.changingThisBreaksApplicationSecurity;
                }
                this.checkNotSafeValue(value, 'ResourceURL');
                throw new Error('unsafe value used in a resource URL context');
            default:
                throw new Error(`Unexpected SecurityContext ${ctx}`);
        }
    }
    checkNotSafeValue(value, expectedType) {
        if (value instanceof SafeValueImpl) {
            throw new Error('Required a safe ' + expectedType + ', got a ' + value.getTypeName());
        }
    }
    sanitizeHtml(value) {
        // TODO(martinprobst): implement.
        return value;
    }
    bypassSecurityTrustHtml(value) { return new SafeHtmlImpl(value); }
    bypassSecurityTrustStyle(value) { return new SafeStyleImpl(value); }
    bypassSecurityTrustScript(value) { return new SafeScriptImpl(value); }
    bypassSecurityTrustUrl(value) { return new SafeUrlImpl(value); }
    bypassSecurityTrustResourceUrl(value) {
        return new SafeResourceUrlImpl(value);
    }
}
DomSanitizationServiceImpl.decorators = [
    { type: Injectable },
];
class SafeValueImpl {
    constructor(changingThisBreaksApplicationSecurity) {
        this.changingThisBreaksApplicationSecurity = changingThisBreaksApplicationSecurity;
        // empty
    }
}
class SafeHtmlImpl extends SafeValueImpl {
    getTypeName() { return 'HTML'; }
}
class SafeStyleImpl extends SafeValueImpl {
    getTypeName() { return 'Style'; }
}
class SafeScriptImpl extends SafeValueImpl {
    getTypeName() { return 'Script'; }
}
class SafeUrlImpl extends SafeValueImpl {
    getTypeName() { return 'URL'; }
}
class SafeResourceUrlImpl extends SafeValueImpl {
    getTypeName() { return 'ResourceURL'; }
}
//# sourceMappingURL=dom_sanitization_service.js.map