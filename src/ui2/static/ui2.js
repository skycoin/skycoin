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

// data-clear="#id" empties a panel. After a send, the form and its password
// are cleared so the same payment cannot go out twice by accident.
(function () {
	'use strict';
	document.addEventListener('click', function (e) {
		var c = e.target.closest && e.target.closest('[data-clear]');
		if (!c) return;
		var p = document.querySelector(c.getAttribute('data-clear'));
		if (p) p.innerHTML = '';
	});
	document.addEventListener('htmx:afterSwap', function (e) {
		if (!e.detail.target || e.detail.target.id !== 'review') return;
		if (!e.detail.target.querySelector('.alert.ok')) return;
		var f = document.getElementById('send');
		if (f) f.reset();
	});
})();

// Transaction notes live in this browser under the Angular wallet's
// "txNotes" key (src/skycoin-web/src/app/utils/tx-notes.ts), so both UIs
// show the same notes.
(function () {
	'use strict';
	var KEY = 'txNotes', MAX = 64;
	function read() {
		try {
			var n = JSON.parse(localStorage.getItem(KEY) || '{}');
			return n && typeof n === 'object' ? n : {};
		} catch (_) { return {}; }
	}
	function fill(root) {
		var notes = read();
		root.querySelectorAll('[data-note]').forEach(function (e) { e.textContent = notes[e.getAttribute('data-note')] || ''; });
		root.querySelectorAll('[data-note-edit]').forEach(function (e) { e.value = notes[e.getAttribute('data-note-edit')] || ''; });
	}
	document.addEventListener('change', function (e) {
		var id = e.target.getAttribute && e.target.getAttribute('data-note-edit');
		if (!id) return;
		var notes = read(), text = e.target.value.trim().substring(0, MAX);
		if (text) notes[id] = text; else delete notes[id];
		try { localStorage.setItem(KEY, JSON.stringify(notes)); } catch (_) { /* private mode */ }
		fill(document);
	});
	document.addEventListener('DOMContentLoaded', function () { fill(document); });
	document.addEventListener('htmx:afterSettle', function (e) { fill(e.detail.elt || document); });
})();
