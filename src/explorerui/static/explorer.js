// Niceties for the server rendered explorer. Every page works without this
// file: it adds copy buttons and shows dates in the reader's time zone.
(function () {
	'use strict';
	var months = ['JAN', 'FEB', 'MAR', 'APR', 'MAY', 'JUN', 'JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'];
	function pad(n) { return (n < 10 ? '0' : '') + n; }
	function zone(d) {
		var m = -d.getTimezoneOffset();
		var s = m < 0 ? '-' : '+';
		m = Math.abs(m);
		return 'GMT' + s + pad(Math.floor(m / 60)) + ':' + pad(m % 60);
	}
	function el(tag, cls, text) {
		var e = document.createElement(tag);
		if (cls) e.className = cls;
		if (text) e.textContent = text;
		return e;
	}

	// Rewrites a UTC <time> in local time, keeping its layout.
	function localize(t) {
		var d = new Date(t.getAttribute('datetime'));
		if (isNaN(d)) return;
		var date = d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate());
		var clock = pad(d.getHours()) + ':' + pad(d.getMinutes());
		t.textContent = '';
		var tz = el('span', 'tz');
		tz.setAttribute('aria-hidden', 'true');
		if (t.getAttribute('data-fmt') === 'stack') {
			var y = el('span', '-gray');
			y.appendChild(tz);
			y.appendChild(document.createTextNode(String(d.getFullYear())));
			t.appendChild(y);
			t.appendChild(el('br'));
			t.appendChild(el('span', '-small-font', months[d.getMonth()] + ' ' + pad(d.getDate())));
			t.appendChild(el('br'));
			t.appendChild(el('span', '-small-font', clock));
		} else {
			t.appendChild(tz);
			t.appendChild(document.createTextNode(date + ' '));
			t.appendChild(el('span', '-gray', clock + ':' + pad(d.getSeconds())));
		}
		t.title = zone(d);
	}

	function enhance(root) {
		var ts = root.querySelectorAll('time[datetime]');
		for (var i = 0; i < ts.length; i++) localize(ts[i]);
		if (navigator.clipboard) {
			var bs = root.querySelectorAll('[data-copy]');
			for (var j = 0; j < bs.length; j++) bs[j].hidden = false;
		}
	}

	document.addEventListener('click', function (e) {
		var b = e.target.closest && e.target.closest('[data-copy]');
		if (!b) return;
		e.preventDefault();
		navigator.clipboard.writeText(b.getAttribute('data-copy')).then(function () {
			b.classList.add('done');
			b.title = document.body.getAttribute('data-copied') || '';
			setTimeout(function () { b.classList.remove('done'); b.title = ''; }, 1500);
		});
	});
	document.addEventListener('htmx:afterSettle', function () { enhance(document); });
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', function () { enhance(document); });
	} else {
		enhance(document);
	}
})();
