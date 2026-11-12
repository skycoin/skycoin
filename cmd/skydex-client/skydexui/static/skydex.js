// Glue for the server rendered SkyDEX client. What the React client keeps in
// localStorage stays there under the same keys, so both clients share it:
// the recent markets, the saved wallet addresses and the trade history.
(function () {
	'use strict';

	var MARKETS = 'exchange:markets';
	var WALLETS = 'exchange:wallets';
	var HISTORY = 'exchange:history';
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

	function outside(el, e) {
		var r = el.getBoundingClientRect();
		return e.clientX < r.left || e.clientX > r.right || e.clientY < r.top || e.clientY > r.bottom;
	}

	// Pages and fragments keep the browser's state in step after each swap.
	function settle(root) {
		// On a phone the tab bar scrolls: keep the open tab in view.
		var at = (!root || root === document.body) && document.querySelector('.tab-link.active');
		if (at) at.scrollIntoView({ block: 'nearest', inline: 'nearest' });
		var h =document.querySelector('header[data-market-pk]');
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
		document.querySelectorAll('time[data-local-time]').forEach(function (el) {
			var d = new Date(el.getAttribute('datetime'));
			if (isNaN(d)) return;
			el.textContent = el.getAttribute('data-local-time') === 'time' ? d.toLocaleTimeString() : d.toLocaleString();
		});
		if (!root || root.querySelector('#create-form')) createForm();
		var hd = document.getElementById('history-data');
		if (hd) {
			try { store(HISTORY, JSON.parse(hd.textContent)); } catch (e) { /* keep what is stored */ }
			hd.remove();
		}
		document.querySelectorAll('[data-autohide]').forEach(function (el) {
			var ms = +el.getAttribute('data-autohide');
			el.removeAttribute('data-autohide');
			setTimeout(function () { el.remove(); }, ms);
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
		var rt = t.closest('[data-remove-trade]');
		if (rt) {
			var id = rt.getAttribute('data-remove-trade');
			store(HISTORY, load(HISTORY, []).filter(function (o) { return o && o.id !== id; }));
			htmx.trigger(document.body, 'history-changed');
			return;
		}
		if (t.closest('[data-clear-trades]')) {
			if (window.confirm('Clear all locally saved trade history on this device?')) {
				store(HISTORY, null);
				htmx.trigger(document.body, 'history-changed');
			}
			return;
		}
		if (t.closest('[data-clear-markets]')) {
			if (window.confirm('Clear all recent markets from this device?')) {
				store(MARKETS, null);
				htmx.trigger(document.body, 'recent-changed');
			}
		}
	});

	// Copy to the clipboard, with a fallback for plain http where
	// navigator.clipboard is missing, as copyToClipboard in src/api.js.
	function copyText(text) {
		if (navigator.clipboard && window.isSecureContext) {
			return navigator.clipboard.writeText(text).then(function () { return true; }, function () { return legacyCopy(text); });
		}
		return Promise.resolve(legacyCopy(text));
	}
	function legacyCopy(text) {
		try {
			var ta = document.createElement('textarea');
			ta.value = text;
			ta.style.position = 'fixed';
			ta.style.top = '-1000px';
			ta.style.opacity = '0';
			document.body.appendChild(ta);
			ta.focus();
			ta.select();
			var ok = document.execCommand('copy');
			document.body.removeChild(ta);
			return ok;
		} catch (e) {
			return false;
		}
	}

	// The QR dialog, as QrButton.jsx.
	var qr = { address: '', uri: '' };
	function qrLib(done) {
		if (window.QRCode) return done();
		var d = document.getElementById('qr-dialog');
		var s = document.createElement('script');
		s.src = d.getAttribute('data-qr-script');
		s.onload = done;
		document.head.appendChild(s);
	}
	function qrDraw() {
		var withAmount = qr.uri && document.getElementById('qr-toggle').checked;
		var value = withAmount ? qr.uri : qr.address;
		var box = document.getElementById('qr-canvas');
		box.textContent = '';
		qrLib(function () {
			new QRCode(box, { text: value, width: 220, height: 220, colorDark: '#000000', colorLight: '#FFFFFF', correctLevel: QRCode.CorrectLevel.M });
		});
		var cu = document.getElementById('qr-copy-uri');
		cu.hidden = !withAmount;
		cu.setAttribute('data-copy', qr.uri);
		document.getElementById('qr-hint').textContent = (withAmount
			? 'Scan with a ' + qr.coin + ' wallet — the address and amount are filled in for you. If your wallet doesn\'t recognise the link, untick the box above for an address-only code.'
			: 'Scan to fill in the address. Enter the amount manually in your wallet.') + (qr.hint ? ' ' + qr.hint : '');
	}
	function qrOpen(b) {
		qr = {
			address: b.getAttribute('data-qr'),
			uri: b.getAttribute('data-qr-link') || '',
			coin: b.getAttribute('data-qr-coin') || '',
			hint: b.getAttribute('data-qr-hint') || ''
		};
		var d = document.getElementById('qr-dialog');
		document.getElementById('qr-title').textContent = b.getAttribute('data-qr-title') || 'Payment';
		var amount = b.getAttribute('data-qr-amount');
		var am = document.getElementById('qr-amount');
		am.hidden = !amount;
		am.querySelector('strong').textContent = amount || '';
		document.getElementById('qr-toggle-row').hidden = !qr.uri;
		document.getElementById('qr-toggle').checked = true;
		document.getElementById('qr-coin').textContent = qr.coin;
		document.getElementById('qr-address').textContent = qr.address;
		document.getElementById('qr-copy-address').setAttribute('data-copy', qr.address);
		qrDraw();
		d.showModal();
	}
	document.addEventListener('change', function (e) {
		if (e.target.id === 'qr-toggle') qrDraw();
		if (e.target.closest('#create-form')) createForm();
	});
	document.addEventListener('input', function (e) {
		if (e.target.closest('#create-form')) createForm();
	});

	// The new sell order form: the coin being sold cannot pay for itself, and
	// the summary follows the inputs, as Market.jsx.
	function createForm() {
		var f = document.getElementById('create-form');
		if (!f) return;
		var sell = f.elements.sell_coin.value;
		var pay = f.elements.payment_currency;
		Array.prototype.forEach.call(pay.options, function (o) {
			o.hidden = o.disabled = o.value === sell;
		});
		if (pay.value === sell || (pay.selectedOptions[0] && pay.selectedOptions[0].disabled)) {
			var first = Array.prototype.filter.call(pay.options, function (o) { return !o.disabled; })[0];
			pay.value = first ? first.value : '';
		}
		f.querySelector('button[type=submit]').disabled = !pay.value;
		var amount = f.elements.amount.value, price = f.elements.price.value;
		var s = document.getElementById('trade-summary');
		s.textContent = '';
		if (amount && price && sell && pay.value) {
			s.append('Sell ', strong(amount + ' ' + sell), ' for ', strong(price + ' ' + pay.value),
				'. After creating this order you have 15 minutes to deposit the ' + sell +
				' (amount + a small commission) to the market escrow; the buyer’s payment goes straight to your wallet.');
		} else {
			s.textContent = 'Set what you sell and what you want for it. Payment can be an external coin (BTC/LTC) or another fibercoin.';
		}
	}
	function strong(t) {
		var el = document.createElement('strong');
		el.textContent = t;
		return el;
	}
	function toggleCreate(open) {
		var panel = document.getElementById('create-panel');
		var t = document.querySelector('[data-toggle-create]');
		if (!panel || !t) return;
		panel.hidden = !open;
		t.classList.toggle('open', open);
	}
	document.body.addEventListener('listing-created', function () {
		var f = document.getElementById('create-form');
		if (f) {
			f.elements.amount.value = '';
			f.elements.price.value = '';
			f.elements.sell_coin.selectedIndex = 0;
			createForm();
		}
		toggleCreate(false);
	});

	// Buying or selling needs saved addresses: send the user to Settings
	// first, as Market.jsx does.
	function needRegister(e) {
		var el = e.target.closest && e.target.closest('[data-need-register]');
		if (el && !registered()) {
			e.preventDefault();
			e.stopPropagation();
			location.href = 'settings';
			return true;
		}
		return false;
	}
	document.addEventListener('submit', needRegister, true);

	document.addEventListener('click', function (e) {
		var t = e.target.closest ? e.target : e.target.parentElement;
		var tc = t.closest('[data-toggle-create]');
		if (tc) {
			e.preventDefault();
			if (needRegister(e)) return;
			toggleCreate(!tc.classList.contains('open'));
			return;
		}
		var c = t.closest('[data-copy]');
		if (c) {
			copyText(c.getAttribute('data-copy')).then(function (ok) {
				if (!ok) return;
				c.textContent = 'Copied ✓';
				setTimeout(function () { c.textContent = c.getAttribute('data-label'); }, 1500);
			});
			return;
		}
		var q = t.closest('[data-qr]');
		if (q) { qrOpen(q); return; }
		if (t.closest('[data-qr-close]') || (t.id === 'qr-dialog' && outside(t, e))) {
			document.getElementById('qr-dialog').close();
			return;
		}
		var d = t.closest('[data-dismiss]');
		if (d) d.closest('.deposit-box').remove();
	});

	settle();
})();
