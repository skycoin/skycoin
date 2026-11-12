// Glue for the server rendered SkyDEX client. What the React client keeps in
// localStorage stays there under the same keys, so both clients share it:
// the recent markets, the saved wallet addresses and the trade history.
(function () {
	'use strict';

	var MARKETS = 'exchange:markets';
	var WALLETS = 'exchange:wallets';
	var MARKETS_MAX = 12;

	function load(key, fallback) {
		try {
			var v = JSON.parse(localStorage.getItem(key));
			return v == null ? fallback : v;
		} catch (e) {
			return fallback;
		}
	}
	function store(key, v) {
		try {
			if (v == null) localStorage.removeItem(key);
			else localStorage.setItem(key, JSON.stringify(v));
		} catch (e) { /* private mode or quota */ }
	}

	function markets() {
		var a = load(MARKETS, []);
		return Array.isArray(a) ? a.filter(function (m) { return m && m.pk; }) : [];
	}
	// Same as addMarketToHistory in src/api.js.
	function addMarket(pk, name) {
		pk = (pk || '').trim();
		if (!pk) return;
		var list = markets();
		var prev = list.filter(function (m) { return m.pk === pk; })[0];
		var rest = list.filter(function (m) { return m.pk !== pk; });
		var nm = (name || '').trim() || (prev && prev.name) || '';
		store(MARKETS, [{ pk: pk, name: nm, ts: Date.now() }].concat(rest).slice(0, MARKETS_MAX));
	}

	function wallets() {
		var w = load(WALLETS, {});
		return w && typeof w === 'object' ? w : {};
	}
	function registered() { return !!wallets().SKY; }

	// Pages and fragments keep the browser's state in step after each swap.
	function settle(root) {
		var h = document.querySelector('header[data-market-pk]');
		if (h && h.getAttribute('data-market-pk')) {
			var n = document.getElementById('market-name');
			addMarket(h.getAttribute('data-market-pk'), n ? n.textContent : h.getAttribute('data-market-name'));
		}
		var b = document.querySelector('[data-register-banner]');
		if (b) b.hidden = registered();
		var w = wallets();
		(root || document).querySelectorAll('input[data-wallet]').forEach(function (el) {
			if (!el.value) el.value = w[el.getAttribute('data-wallet')] || '';
		});
	}

	// Elements with data-local post that localStorage key as "saved", so the
	// server can render the browser's own list.
	document.addEventListener('htmx:configRequest', function (e) {
		var key = e.detail.elt.getAttribute && e.detail.elt.getAttribute('data-local');
		if (key) e.detail.parameters.saved = JSON.stringify(load(key, []));
	});

	document.addEventListener('htmx:afterSettle', function (e) { settle(e.detail.elt); });

	// The market accepted the addresses: remember them, as Settings.jsx does.
	document.body.addEventListener('wallets-saved', function () {
		var w = {};
		document.querySelectorAll('#wallets-form input[data-wallet]').forEach(function (el) {
			var v = el.value.trim();
			if (v) w[el.getAttribute('data-wallet')] = v;
		});
		store(WALLETS, w);
		settle();
	});

	document.addEventListener('click', function (e) {
		var t = e.target.closest ? e.target : e.target.parentElement;
		var del = t.closest('[data-remove-market]');
		if (del) {
			var pk = del.getAttribute('data-remove-market');
			store(MARKETS, markets().filter(function (m) { return m.pk !== pk; }));
			htmx.trigger(document.body, 'recent-changed');
			return;
		}
		if (t.closest('[data-clear-markets]')) {
			if (window.confirm('Clear all recent markets from this device?')) {
				store(MARKETS, null);
				htmx.trigger(document.body, 'recent-changed');
			}
		}
	});

	settle();
})();
