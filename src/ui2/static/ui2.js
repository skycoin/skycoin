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
			.then(function (r) { return r.ok ? r.json() : {}; })
			// A host that hands out no token, or cannot, gets the request without
			// one; a host that checks tokens then refuses it.
			.catch(function () { return {}; })
			.then(function (j) { token = (j && j.csrf_token) || ''; e.detail.issueRequest(true); });
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
		var cc = e.target.closest && e.target.closest('[data-clear-closest]');
		if (cc) {
			var box = cc.closest('form').querySelector(cc.getAttribute('data-clear-closest'));
			if (box) box.innerHTML = '';
			return;
		}
		var c = e.target.closest && e.target.closest('[data-clear]');
		if (!c) return;
		var p = document.querySelector(c.getAttribute('data-clear'));
		if (p) p.innerHTML = '';
	});
	document.addEventListener('htmx:afterSwap', function (e) {
		if (!e.detail.target || e.detail.target.id !== 'review') return;
		if (!e.detail.target.querySelector('.alert.ok')) return;
		saveNote(e.detail.target.querySelector('[data-sent-txid]'), document.getElementById('send-note'));
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

// "Send all" fills the only destination with every confirmed coin.
document.addEventListener('click', function (e) {
	var b = e.target.closest && e.target.closest('[data-fill-all]');
	if (!b) return;
	var rows = document.querySelectorAll('#send-rows [name=coins]');
	if (rows.length) rows[0].value = b.getAttribute('data-fill-all');
});

// saveNote keeps the note typed on the send form for the transaction the
// node just took, under the Angular wallet's txNotes key.
function saveNote(sent, input) {
	if (!sent || !input || !input.value.trim()) return;
	try {
		var notes = JSON.parse(localStorage.getItem('txNotes') || '{}') || {};
		notes[sent.getAttribute('data-sent-txid')] = input.value.trim().substring(0, 64);
		localStorage.setItem('txNotes', JSON.stringify(notes));
	} catch (_) { /* private mode */ }
}

// The QR dialog. The content follows the Angular wallet's QR dialog:
// coin:ADDRESS, then amount, hours and message when asked for.
(function () {
	'use strict';
	var dlg = document.getElementById('qr-dialog');
	if (!dlg) return;
	var address = '', coin = 'skycoin', loading = null;
	function $(id) { return document.getElementById(id); }
	function lib() {
		if (window.QRCode) return Promise.resolve();
		if (!loading) {
			loading = new Promise(function (resolve, reject) {
				var s = document.createElement('script');
				s.src = dlg.getAttribute('data-qr-script');
				s.onload = resolve;
				s.onerror = reject;
				document.head.appendChild(s);
			});
		}
		return loading;
	}
	function content() {
		var text = coin + ':' + address, sep = '?', bad = false;
		var c = $('qr-coins').value.trim(), h = $('qr-hours').value.trim(), m = $('qr-note').value.trim();
		if (c) {
			if (String(parseFloat(c)) === c && parseFloat(c) > 0) { text += sep + 'amount=' + c; sep = '&'; } else bad = true;
		}
		if (h) {
			if (String(parseInt(h, 10)) === h && parseInt(h, 10) > 0) { text += sep + 'hours=' + h; sep = '&'; } else bad = true;
		}
		if (m) text += sep + 'message=' + encodeURIComponent(m);
		$('qr-invalid').hidden = !bad;
		return text;
	}
	function draw() {
		var text = content();
		$('qr-data').textContent = text;
		$('qr-data').setAttribute('data-copy', text);
		lib().then(function () {
			var box = $('qr-box');
			box.innerHTML = '';
			new window.QRCode(box, { text: text, width: 240, height: 240, colorDark: '#000000', colorLight: '#ffffff', correctLevel: window.QRCode.CorrectLevel.M });
		});
	}
	var timer = 0;
	dlg.addEventListener('input', function () { clearTimeout(timer); timer = setTimeout(draw, 300); });
	document.addEventListener('click', function (e) {
		var b = e.target.closest && e.target.closest('[data-qr]');
		if (!b) return;
		address = b.getAttribute('data-qr');
		coin = b.getAttribute('data-qr-coin') || 'skycoin';
		$('qr-addr').textContent = address;
		$('qr-addr').setAttribute('data-copy', address);
		$('qr-coins').value = $('qr-hours').value = $('qr-note').value = '';
		draw();
		dlg.showModal();
	});
})();

// A payment link fills the first destination, as the Angular wallet's
// parsePaymentLink reads it (src/skycoin-web/src/app/utils/qr.ts).
document.addEventListener('click', function (e) {
	var b = e.target.closest && e.target.closest('[data-fill-link]');
	if (!b) return;
	var box = document.getElementById('pay-link');
	var text = document.getElementById('pay-link-text').value.trim();
	var m = /^([a-z][a-z0-9+.-]*):(?:\/\/)?([^?]*)(?:\?(.*))?$/i.exec(text);
	var q = m && m[3] ? new URLSearchParams(m[3]) : new URLSearchParams('');
	var addr = m ? m[2].split('@')[0].split('/')[0] : '';
	var amount = q.get('amount'), hours = q.get('hours');
	if (!addr || (amount !== null && !/^\d*\.?\d+$/.test(amount)) || (hours !== null && !/^\d+$/.test(hours))) {
		toastMsg(box.getAttribute('data-msg-invalid'), true);
		return;
	}
	var f = document.getElementById('send-rows').closest('form');
	f.querySelector('#send-rows [name=to]').value = addr;
	if (amount !== null) f.querySelector('#send-rows [name=coins]').value = amount;
	if (hours !== null) {
		f.querySelector('#send-rows [name=hours]').value = hours;
		var manual = f.querySelector('[name=hours_mode][value=manual]');
		if (manual) manual.checked = true;
	}
	toastMsg(box.getAttribute('data-msg-done'), false);
});

function toastMsg(msg, isErr) {
	var t = document.getElementById('toast');
	if (!t) return;
	t.textContent = msg;
	t.className = 'show' + (isErr ? ' err' : '');
	setTimeout(function () { t.className = ''; }, 3000);
}

// History lists keep their last good copy in this browser. When the node
// cannot answer, the copy is shown with the time it was saved, as the
// Angular history does.
(function () {
	'use strict';
	var PREFIX = 'ui2History:';
	function sync(root) {
		var lists = root.matches && root.matches('[data-history-cache]') ? [root] : root.querySelectorAll('[data-history-cache]');
		Array.prototype.forEach.call(lists, function (el) {
			var key = el.getAttribute('data-history-cache');
			if (!key) return;
			try {
				if (!el.querySelector('.error')) {
					localStorage.setItem(PREFIX + key, JSON.stringify({ at: Date.now(), html: el.innerHTML }));
					return;
				}
				var saved = JSON.parse(localStorage.getItem(PREFIX + key) || 'null');
				if (!saved) return;
				var note = document.createElement('p');
				note.className = 'alert warn cached-note';
				note.textContent = (document.body.getAttribute('data-cached') || '').replace('{{ time }}', new Date(saved.at).toLocaleString());
				el.innerHTML = saved.html;
				el.insertBefore(note, el.firstChild);
			} catch (_) { /* storage blocked */ }
		});
	}
	document.addEventListener('DOMContentLoaded', function () { sync(document); });
	document.addEventListener('htmx:afterSettle', function () { sync(document); });
})();

// fillPayment puts a scanned or pasted address and amount in the send
// form's first destination.
function fillPayment(addr, amount, hours) {
	var rows = document.getElementById('send-rows');
	if (!rows) return false;
	var f = rows.closest('form');
	rows.querySelector('[name=to]').value = addr;
	if (amount) rows.querySelector('[name=coins]').value = amount;
	if (hours) {
		rows.querySelector('[name=hours]').value = hours;
		var manual = f.querySelector('[name=hours_mode][value=manual]');
		if (manual) manual.checked = true;
	}
	return true;
}

// The QR scan dialog. A scanned text is a bare address or a payment URI
// such as skycoin:ADDR?amount=1, read as the Angular wallet's
// parseScannedPayment does (src/skycoin-web/src/app/utils/qr.ts).
(function () {
	'use strict';
	var dlg = document.getElementById('scan-dialog');
	if (!dlg) return;
	var video = document.getElementById('scan-video'), status = document.getElementById('scan-status');
	var canvas = document.createElement('canvas'), stream = null, timer = 0, loading = null;
	function lib() {
		if (window.jsQR) return Promise.resolve();
		if (!loading) {
			loading = new Promise(function (resolve, reject) {
				var s = document.createElement('script');
				s.src = dlg.getAttribute('data-jsqr');
				s.onload = resolve;
				s.onerror = reject;
				document.head.appendChild(s);
			});
		}
		return loading;
	}
	function decode(src, w, h) {
		var scale = Math.min(1, 640 / Math.max(w, h));
		canvas.width = Math.max(1, Math.round(w * scale));
		canvas.height = Math.max(1, Math.round(h * scale));
		var ctx = canvas.getContext('2d', { willReadFrequently: true });
		ctx.drawImage(src, 0, 0, canvas.width, canvas.height);
		var img = ctx.getImageData(0, 0, canvas.width, canvas.height);
		var code = window.jsQR(img.data, img.width, img.height);
		return code ? code.data : null;
	}
	function parsed(text) {
		var rest = text.trim(), scheme = /^([a-z][a-z0-9+.-]*):/i.exec(rest);
		if (scheme) rest = rest.substring(scheme[0].length).replace(/^\/\//, '');
		var parts = rest.split('?'), addr = parts[0].split('@')[0].split('/')[0], amount = null, hours = null;
		if (parts[1] && !(scheme && scheme[1].toLowerCase() === 'ethereum')) {
			var q = new URLSearchParams(parts[1]);
			amount = q.get('amount');
			if (amount && !/^\d*\.?\d+$/.test(amount)) amount = null;
			hours = q.get('hours');
			if (hours && !/^\d+$/.test(hours)) hours = null;
		}
		return { addr: addr, amount: amount, hours: hours };
	}
	function stop() {
		clearInterval(timer);
		if (stream) stream.getTracks().forEach(function (t) { t.stop(); });
		stream = null;
		video.hidden = true;
	}
	function found(text) {
		var p = parsed(text);
		if (!p.addr || !fillPayment(p.addr, p.amount, p.hours)) return false;
		stop();
		dlg.close();
		if (typeof toastMsg === 'function') toastMsg(dlg.getAttribute('data-msg-done'), false);
		return true;
	}
	dlg.addEventListener('close', stop);
	document.getElementById('scan-file').addEventListener('change', function (e) {
		var file = e.target.files && e.target.files[0];
		if (!file) return;
		// A data URL, since the page's CSP allows data: images and not blob: ones.
		var img = new Image(), reader = new FileReader();
		var none = function () { status.textContent = dlg.getAttribute('data-msg-no-code'); };
		img.onerror = none;
		img.onload = function () {
			lib().then(function () {
				var text = decode(img, img.naturalWidth, img.naturalHeight);
				if (!text || !found(text)) none();
			}, none);
		};
		reader.onload = function () { img.src = reader.result; };
		reader.onerror = none;
		reader.readAsDataURL(file);
		e.target.value = '';
	});
	document.addEventListener('click', function (e) {
		if (!e.target.closest || !e.target.closest('[data-qr-scan]')) return;
		status.textContent = '';
		dlg.showModal();
		lib();
		if (!window.isSecureContext || !navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
			status.textContent = dlg.getAttribute('data-msg-insecure');
			return;
		}
		navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } }).then(function (s) {
			stream = s;
			video.srcObject = s;
			video.hidden = false;
			return video.play();
		}).then(function () {
			timer = setInterval(function () {
				if (!window.jsQR || !video.videoWidth) return;
				var text = decode(video, video.videoWidth, video.videoHeight);
				if (text) found(text);
			}, 250);
		}).catch(function () {
			stop();
			status.textContent = dlg.getAttribute('data-msg-no-camera');
		});
	});
})();

// The coin picker keeps the chosen coin in a cookie the server reads, and
// opens that coin's wallets.
document.addEventListener('change', function (e) {
	if (!e.target.matches || !e.target.matches('[data-coin]')) return;
	document.cookie = 'ui2-coin=' + e.target.value + ';path=/;max-age=31536000;samesite=strict';
	window.location.href = document.querySelector('base').href;
});
