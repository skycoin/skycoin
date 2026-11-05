// Wallets kept in this browser, in the Angular wallet's localStorage format
// (src/skycoin-web/src/app/services/wallet/backends/browser-wallet-backend.ts),
// so a wallet made in one UI opens in the other. Seeds, passwords and keys
// stay in this script and the wasm cipher. The server is sent public data
// only: addresses, the account xpub and labels.
(function () {
	'use strict';

	var page = document.querySelector('[data-bw-page]');
	if (!page) return;
	var mode = page.getAttribute('data-bw-page');
	var root = page.getAttribute('data-bw-root') || '../';
	var WALLETS = 'wallets', LOCKED = 'lockedWallets';
	var appPassword = null;

	function $(id) { return document.getElementById(id); }
	function msg(name) { return page.getAttribute('data-msg-' + name) || name; }
	function normalize(s) { return (s || '').trim().split(/\s+/).filter(Boolean).join(' '); }

	// Encryption as src/skycoin-web/src/app/services/encryption.service.ts:
	// PBKDF2-SHA256 with 100000 rounds to an AES-GCM key, packed as base64 of
	// salt (16) + iv (12) + ciphertext.
	var enc = new TextEncoder(), dec = new TextDecoder();
	function deriveKey(password, salt) {
		return crypto.subtle.importKey('raw', enc.encode(password), 'PBKDF2', false, ['deriveKey']).then(function (k) {
			return crypto.subtle.deriveKey({ name: 'PBKDF2', salt: salt, iterations: 100000, hash: 'SHA-256' },
				k, { name: 'AES-GCM', length: 256 }, false, ['encrypt', 'decrypt']);
		});
	}
	function toB64(bytes) {
		var s = '';
		for (var i = 0; i < bytes.length; i++) s += String.fromCharCode(bytes[i]);
		return btoa(s);
	}
	function encrypt(text, password) {
		var salt = crypto.getRandomValues(new Uint8Array(16));
		var iv = crypto.getRandomValues(new Uint8Array(12));
		return deriveKey(password, salt).then(function (key) {
			return crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, key, enc.encode(text));
		}).then(function (ct) {
			var packed = new Uint8Array(28 + ct.byteLength);
			packed.set(salt, 0);
			packed.set(iv, 16);
			packed.set(new Uint8Array(ct), 28);
			return toB64(packed);
		});
	}
	function decrypt(b64, password) {
		var packed = Uint8Array.from(atob(b64), function (c) { return c.charCodeAt(0); });
		return deriveKey(password, packed.slice(0, 16)).then(function (key) {
			return crypto.subtle.decrypt({ name: 'AES-GCM', iv: packed.slice(16, 28) }, key, packed.slice(28));
		}).then(function (pt) { return dec.decode(pt); });
	}

	// Storage. The whole list may be encrypted under an app password, which
	// is asked for here and kept in memory for this page only.
	function readStored() {
		var locked = null;
		try { locked = localStorage.getItem(LOCKED); } catch (_) { return Promise.reject(new Error(msg('error'))); }
		if (locked) {
			if (appPassword === null) return Promise.resolve(null);
			return decrypt(locked, appPassword).then(JSON.parse);
		}
		var plain = localStorage.getItem(WALLETS);
		return Promise.resolve(plain ? JSON.parse(plain) : []);
	}
	function writeStored(list) {
		var s = JSON.stringify(list);
		if (localStorage.getItem(LOCKED)) {
			return encrypt(s, appPassword).then(function (c) {
				localStorage.setItem(LOCKED, c);
				localStorage.removeItem(WALLETS);
			});
		}
		localStorage.setItem(WALLETS, s);
		return Promise.resolve();
	}
	// withList runs f on the stored list, first asking for the app password
	// when the list is locked.
	function withList(f) {
		return readStored().then(function (list) {
			if (list !== null) return f(list);
			return new Promise(function (resolve, reject) {
				var form = $('bw-unlock');
				form.hidden = false;
				form.addEventListener('submit', function (e) {
					e.preventDefault();
					var pw = $('bw-app-password').value;
					decrypt(localStorage.getItem(LOCKED), pw).then(function (s) {
						appPassword = pw;
						$('bw-app-password').value = '';
						form.hidden = true;
						resolve(f(JSON.parse(s)));
					}, function () {
						$('bw-unlock-error').textContent = 'Incorrect password';
					});
				});
			});
		});
	}

	function isBip44(w) { return w.walletType === 'bip44' && !w.filename && w.accounts && w.accounts.length > 0; }
	function ours(w) { return !w.isHardware && !w.coinId; }

	// summary is the public part of a stored wallet, as ui2.BrowserWallet.
	function summary(w) {
		var s = { label: w.label || '', encrypted: !!w.encryptedSeed, addresses: [] };
		if (isBip44(w)) {
			var a = w.accounts[0];
			s.type = 'bip44';
			s.xpub = a.accountXpubKey || '';
			a.externalAddresses.forEach(function (e, i) { s.addresses.push({ address: e.address, chain: 'external', index: i }); });
			a.changeAddresses.forEach(function (e, i) { s.addresses.push({ address: e.address, chain: 'change', index: i }); });
		} else {
			s.type = w.walletType || 'deterministic';
			w.addresses.forEach(function (e, i) { s.addresses.push({ address: e.address, index: i }); });
		}
		s.id = s.addresses.length ? s.addresses[0].address : '';
		return s;
	}
	function firstAddress(w) {
		if (isBip44(w)) return w.accounts[0].externalAddresses[0] && w.accounts[0].externalAddresses[0].address;
		return w.addresses[0] && w.addresses[0].address;
	}

	function ready(value) {
		$('bw-data').value = JSON.stringify(value);
		document.body.dispatchEvent(new CustomEvent('bw-ready'));
	}

	// The wasm cipher, the same build the Angular wallet loads.
	var cipherReady = null;
	function cipher() {
		if (!cipherReady) {
			cipherReady = (window.SkycoinCipher ? Promise.resolve() : (function () {
				var go = new window.Go();
				return WebAssembly.instantiateStreaming(fetch(root + 'assets/scripts/skycoin-lite.wasm'), go.importObject)
					.then(function (r) { go.run(r.instance); });
			})()).then(function () { return window.SkycoinCipher; });
		}
		return cipherReady;
	}
	function call(c, name, args) {
		var r = c[name].apply(c, args);
		if (r === null || r === undefined) throw new Error('the cipher could not ' + name);
		if (r.error) throw new Error(r.error);
		return r;
	}
	function asciiToHex(s) {
		var out = '';
		for (var i = 0; i < s.length; i++) {
			var h = s.charCodeAt(i).toString(16);
			out += h.length === 1 ? '0' + h : h;
		}
		return out;
	}

	// BIP39 with the word list the Angular wallet ships.
	var wordsReady = null;
	function wordList() {
		if (!wordsReady) {
			wordsReady = fetch(root + 'assets/bip39-word-list.json').then(function (r) { return r.json(); })
				.then(function (j) { return j.list; });
		}
		return wordsReady;
	}
	function bitsOf(bytes) {
		var s = '';
		bytes.forEach(function (b) { s += ('00000000' + b.toString(2)).slice(-8); });
		return s;
	}
	function newMnemonic(words) {
		var entropy = crypto.getRandomValues(new Uint8Array(words === 24 ? 32 : 16));
		return Promise.all([wordList(), crypto.subtle.digest('SHA-256', entropy)]).then(function (r) {
			var bits = bitsOf(entropy) + bitsOf(new Uint8Array(r[1])).slice(0, entropy.length / 4);
			var out = [];
			for (var i = 0; i < bits.length; i += 11) out.push(r[0][parseInt(bits.slice(i, i + 11), 2)]);
			return out.join(' ');
		});
	}
	function validMnemonic(m) {
		var ws = m.split(' ');
		if ([12, 15, 18, 21, 24].indexOf(ws.length) < 0) return Promise.resolve(false);
		return wordList().then(function (list) {
			var bits = '';
			for (var i = 0; i < ws.length; i++) {
				var idx = list.indexOf(ws[i]);
				if (idx < 0) return false;
				bits += ('00000000000' + idx.toString(2)).slice(-11);
			}
			var cs = ws.length / 3, ent = bits.slice(0, bits.length - cs);
			var bytes = new Uint8Array(ent.length / 8);
			for (var j = 0; j < bytes.length; j++) bytes[j] = parseInt(ent.slice(j * 8, j * 8 + 8), 2);
			return crypto.subtle.digest('SHA-256', bytes).then(function (h) {
				return bitsOf(new Uint8Array(h)).slice(0, cs) === bits.slice(bits.length - cs);
			});
		});
	}

	// used is the addresses in addrs that appear in any transaction.
	function used(addrs) {
		if (!addrs.length) return Promise.resolve({});
		return fetch(root + 'api/v1/transactions?verbose=1&addrs=' + encodeURIComponent(addrs.join(',')), { credentials: 'same-origin' })
			.then(function (r) { if (!r.ok) throw new Error(r.status + ' ' + r.statusText); return r.json(); })
			.then(function (txns) {
				var want = {}, seen = {};
				addrs.forEach(function (a) { want[a] = true; });
				(txns || []).forEach(function (t) {
					(t.txn.inputs || []).forEach(function (i) { if (want[i.owner]) seen[i.owner] = true; });
					(t.txn.outputs || []).forEach(function (o) { if (want[o.dst]) seen[o.dst] = true; });
				});
				return seen;
			});
	}
	// scan keeps the addresses up to the last used one, deriving batches of
	// 20 until a batch has none in use, the BIP44 gap limit.
	function scan(derive, keepAtLeast) {
		var all = [], last = -1;
		function step() {
			var batch = derive(all.length, 20);
			return used(batch.map(function (a) { return a.address; })).then(function (seen) {
				var any = false;
				batch.forEach(function (a) {
					all.push(a);
					if (seen[a.address]) { last = all.length - 1; any = true; }
				});
				return any ? step() : all.slice(0, Math.max(last + 1, keepAtLeast));
			});
		}
		return step();
	}

	function setError(text) { var e = $('bw-error'); if (e) e.textContent = text || ''; }
	function setStatus(text) { var e = $('bw-status'); if (e) e.textContent = text || ''; }

	function bip44Entry(label, xpub, external, change, encSeed, encPassphrase) {
		var ext = external.map(function (a) { return { address: a.address }; });
		var chg = change.map(function (a) { return { address: a.address }; });
		var w = {
			coinId: 0, needSeedConfirmation: false, label: label, addresses: ext.concat(chg), walletType: 'bip44',
			accounts: [{ name: 'Account 0', index: 0, accountXpubKey: xpub, externalAddresses: ext, changeAddresses: chg }],
			encryptedSeed: encSeed,
		};
		if (encPassphrase) w.encryptedSeedPassphrase = encPassphrase;
		return w;
	}

	// save adds w to the stored list unless a wallet with its first address
	// is there, then opens it.
	function save(w) {
		return withList(function (list) {
			var first = firstAddress(w);
			if (list.some(function (x) { return ours(x) && firstAddress(x) === first; })) throw new Error(msg('exists'));
			list.push(w);
			return writeStored(list).then(function () { window.location.href = 'bw/wallet?id=' + encodeURIComponent(first); });
		});
	}

	function createFlow(load) {
		var generated = null;
		var gen = $('bw-generate');
		if (gen) gen.addEventListener('click', function () {
			newMnemonic(parseInt($('bw-words').value, 10)).then(function (m) {
				generated = m;
				var ol = $('bw-seed-words');
				ol.textContent = '';
				m.split(' ').forEach(function (w) {
					var li = document.createElement('li');
					li.className = 'mono';
					li.textContent = w;
					ol.appendChild(li);
				});
				$('bw-seed-step').hidden = false;
			});
		});
		$('bw-create').addEventListener('submit', function (e) {
			e.preventDefault();
			setError('');
			var label = $('bw-label').value.trim();
			var pw = $('bw-password').value, pw2 = $('bw-password2').value;
			var seed = load ? normalize($('bw-seed').value) : generated;
			var type = load ? $('bw-type').value : 'bip44';
			var passphrase = load ? $('bw-passphrase').value : '';
			if (!label) return setError(msg('label'));
			if (!seed) return setError(msg('seed'));
			if (!load && normalize($('bw-seed-confirm').value) !== seed) return setError(msg('seed-confirm'));
			if (!pw) return setError(msg('password'));
			if (pw !== pw2) return setError(msg('password-confirm'));
			var btn = $('bw-submit');
			btn.disabled = true;
			validMnemonic(seed).then(function (ok) {
				if (!ok && (type === 'bip44' || !$('bw-nonstandard').checked)) throw new Error(msg('not-bip39'));
				return cipher();
			}).then(function (c) {
				setStatus(msg('scanning'));
				if (type === 'bip44') {
					var xpub = call(c, 'bip44AccountXPub', [seed, passphrase]);
					var derive = function (chain) {
						return function (start, n) { return call(c, 'bip44XPubAddresses', [xpub, chain, start, n]); };
					};
					var first = derive(0)(0, 1);
					var scanned = load
						? Promise.all([scan(derive(0), 1), scan(derive(1), 0)])
						: Promise.resolve([first, []]);
					return scanned.then(function (r) {
						return Promise.all([encrypt(seed, pw), passphrase ? encrypt(passphrase, pw) : null]).then(function (secrets) {
							return bip44Entry(label, xpub, r[0], r[1], secrets[0], secrets[1]);
						});
					});
				}
				// A deterministic wallet: each address's seed is the next
				// seed of the one before, from the hex of the seed text.
				var chain = [];
				var next = asciiToHex(seed);
				var derive = function (start, n) {
					while (chain.length < start + n) {
						var a = call(c, 'generateAddress', [next]);
						chain.push({ address: a.address, nextSeed: a.nextSeed });
						next = a.nextSeed;
					}
					return chain.slice(start, start + n);
				};
				return scan(derive, 1).then(function (kept) {
					var after = kept[kept.length - 1].nextSeed;
					return Promise.all([encrypt(seed, pw), encrypt(after, pw)]).then(function (secrets) {
						return {
							coinId: 0, needSeedConfirmation: false, label: label,
							addresses: kept.map(function (a) { return { address: a.address }; }),
							encryptedSeed: secrets[0], encryptedNextSeed: secrets[1],
						};
					});
				});
			}).then(save).catch(function (err) {
				setStatus('');
				setError(err && err.message ? err.message : String(err));
				btn.disabled = false;
			}).then(function () {
				seed = null; pw = null; pw2 = null; passphrase = null;
			});
		});
	}

	document.addEventListener('DOMContentLoaded', function () {
		var fail = function (err) {
			var e = $('bw-page-error') || $('bw-error');
			if (e) e.textContent = err && err.message ? err.message : msg('error');
		};
		if (mode === 'list') {
			withList(function (list) { ready(list.filter(ours).map(summary)); }).catch(fail);
		} else if (mode === 'wallet') {
			var id = page.getAttribute('data-bw-id');
			withList(function (list) {
				var w = list.filter(ours).find(function (x) { return firstAddress(x) === id; });
				if (!w) throw new Error(msg('missing'));
				ready(summary(w));
			}).catch(fail);
		} else if (mode === 'new' || mode === 'load') {
			createFlow(mode === 'load');
		}
	});
})();
