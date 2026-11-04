// Glue for the htmx pages: a fresh CSRF token per write, toasts, copy to
// clipboard, readable errors, the theme, the language and the phone menu.
(function () {
	'use strict';
	var token = '';
	var toastTimer = 0;

	function toast(msg, isErr) {
		var t = document.getElementById('toast');
		if (!t) return;
		t.textContent = msg;
		t.className = 'show' + (isErr ? ' err' : '');
		clearTimeout(toastTimer);
		toastTimer = setTimeout(function () { t.className = ''; }, 3000);
	}
	function store(k, v) {
		try { if (v === null) localStorage.removeItem(k); else localStorage.setItem(k, v); } catch (_) { /* private mode */ }
	}
	function load(k) {
		try { return localStorage.getItem(k); } catch (_) { return null; }
	}

	// The node checks X-CSRF-Token on every write and a token lives 30s, so
	// fetch one right before the request (src/api/csrf.go).
	document.addEventListener('htmx:confirm', function (e) {
		if (e.detail.verb === 'get') return;
		e.preventDefault();
		if (e.detail.question && !window.confirm(e.detail.question)) return;
		fetch(document.body.getAttribute('data-csrf'), { credentials: 'same-origin' })
			.then(function (r) { return r.status === 404 ? {} : r.json(); })
			.then(function (j) { token = j.csrf_token || ''; e.detail.issueRequest(true); })
			.catch(function (err) { toast('csrf: ' + err, true); });
	});
	document.addEventListener('htmx:configRequest', function (e) {
		if (e.detail.verb !== 'get' && token) e.detail.headers['X-CSRF-Token'] = token;
	});

	document.addEventListener('htmx:responseError', function (e) {
		var x = e.detail.xhr;
		toast(x.status + ' ' + (x.responseText || x.statusText).slice(0, 200), true);
	});
	document.addEventListener('htmx:sendError', function () { toast(document.body.getAttribute('data-unreachable'), true); });

	document.addEventListener('click', function (e) {
		var k = e.target.closest && e.target.closest('[data-copy]');
		if (k) {
			e.preventDefault();
			if (navigator.clipboard) {
				navigator.clipboard.writeText(k.getAttribute('data-copy')).then(function () { toast('✓'); });
			}
			return;
		}
		if (e.target.closest && e.target.closest('[data-theme-toggle]')) {
			var cur = document.documentElement.getAttribute('data-theme');
			var next = cur === 'light' ? 'dark' : cur === 'dark' ? null : 'light';
			if (next) document.documentElement.setAttribute('data-theme', next);
			else document.documentElement.removeAttribute('data-theme');
			store('ui2-theme', next);
			return;
		}
		var m = e.target.closest && e.target.closest('[data-menu-toggle]');
		if (m) {
			var top = m.closest('.top');
			var open = top.classList.toggle('open');
			m.setAttribute('aria-expanded', open ? 'true' : 'false');
		}
	});

	// Language: the server reads a cookie, the Angular wallet reads
	// localStorage "lang".
	function setLang(code) {
		document.cookie = 'ui2-lang=' + code + ';path=/;max-age=31536000;samesite=strict';
		store('lang', code);
	}
	document.addEventListener('change', function (e) {
		if (!e.target.matches || !e.target.matches('[data-lang]')) return;
		setLang(e.target.value);
		window.location.reload();
	});
	if (document.cookie.indexOf('ui2-lang=') < 0) {
		var ng = load('lang');
		if (ng && ng !== document.documentElement.lang) { setLang(ng); window.location.reload(); }
		else if (document.documentElement.lang) setLang(document.documentElement.lang);
	}
})();
