// Applies a chosen light or dark theme before the page paints. Without a
// choice the stylesheet follows the system. A file, not inline, for the CSP.
try {
	var t = localStorage.getItem('ui2-theme');
	if (t === 'light' || t === 'dark') document.documentElement.setAttribute('data-theme', t);
} catch (e) { /* private mode */ }
