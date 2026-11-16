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
	// The keys need WebCrypto, which a page outside a secure context lacks.
	if (!window.isSecureContext || !window.crypto || !window.crypto.subtle) {
		var warn = document.createElement('div');
		warn.className = 'alert warn';
		warn.setAttribute('role', 'alert');
		warn.textContent = document.body.getAttribute('data-insecure');
		page.insertBefore(warn, page.firstChild);
	}
	var root = page.getAttribute('data-bw-root') || '';
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
			return decrypt(locked, appPassword).then(function (s) { return JSON.parse(s).concat(readTemp()); });
		}
		var plain = localStorage.getItem(WALLETS);
		return Promise.resolve((plain ? JSON.parse(plain) : []).concat(readTemp()));
	}
	// Temporary wallets live in this tab's sessionStorage, still encrypted,
	// and are gone when the tab closes. The Angular wallet keeps them in
	// memory only, so it never sees them.
	var TEMP = 'ui2TempWallets';
	function readTemp() {
		try { return JSON.parse(sessionStorage.getItem(TEMP) || '[]'); } catch (_) { return []; }
	}
	function writeStored(list) {
		try { sessionStorage.setItem(TEMP, JSON.stringify(list.filter(function (w) { return w.temporary; }))); } catch (_) { /* private mode */ }
		var s = JSON.stringify(list.filter(function (w) { return !w.temporary; }));
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
	// A Skywallet's wallet holds addresses only. static/hw.js has the device sign.
	function isHw(w) { return !!w.isHardware && !w.coinId; }
	function shown(w) { return ours(w) || isHw(w); }

	// summary is the public part of a stored wallet, as ui2.BrowserWallet.
	function summary(w) {
		var s = { label: w.label || '', encrypted: !!w.encryptedSeed, addresses: [] };
		if (w.temporary) s.temporary = true;
		if (isBip44(w)) {
			var a = w.accounts[0];
			s.type = 'bip44';
			s.xpub = a.accountXpubKey || '';
			a.externalAddresses.forEach(function (e, i) { s.addresses.push({ address: e.address, chain: 'external', index: i }); });
			a.changeAddresses.forEach(function (e, i) { s.addresses.push({ address: e.address, chain: 'change', index: i }); });
		} else {
			s.type = isHw(w) ? 'skywallet' : w.walletType || 'deterministic';
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
	// initialState is what the first address of a deterministic or
	// post-quantum wallet with seed comes from, in the Angular wallet's form:
	// the seed's hex, or pq:<hex>:<index> for a post-quantum wallet.
	function initialState(seed, type) {
		return type === 'pq' ? 'pq:' + asciiToHex(seed) + ':0' : asciiToHex(seed);
	}
	// step derives the address of state, its key and the state after it.
	function step(c, state) {
		var pq = /^pq:([0-9a-f]*):(\d+)$/.exec(state);
		if (!pq) return call(c, 'generateAddress', [state]);
		var i = parseInt(pq[2], 10);
		var a = call(c, 'generatePQAddress', [pq[1], i]);
		return { address: a.address, secret: a.secret, nextSeed: 'pq:' + pq[1] + ':' + (i + 1) };
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
			var type = $('bw-type').value;
			var passphrase = type === 'bip44' ? $('bw-passphrase').value : '';
			var temp = load && $('bw-temp').checked;
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
				// A deterministic or post-quantum wallet: each address comes
				// from the state the one before it leaves.
				var chain = [];
				var next = initialState(seed, type);
				var derive = function (start, n) {
					while (chain.length < start + n) {
						var a = step(c, next);
						chain.push({ address: a.address, nextSeed: a.nextSeed });
						next = a.nextSeed;
					}
					return chain.slice(start, start + n);
				};
				return scan(derive, 1).then(function (kept) {
					var after = kept[kept.length - 1].nextSeed;
					return Promise.all([encrypt(seed, pw), encrypt(after, pw)]).then(function (secrets) {
						return Object.assign({
							coinId: 0, needSeedConfirmation: false, label: label,
							addresses: kept.map(function (a) { return { address: a.address }; }),
							encryptedSeed: secrets[0], encryptedNextSeed: secrets[1],
						}, type === 'pq' ? { walletType: 'pq' } : {});
					});
				});
			}).then(function (w) {
				if (temp) w.temporary = true;
				return save(w);
			}).catch(function (err) {
				setStatus('');
				setError(err && err.message ? err.message : String(err));
				btn.disabled = false;
			}).then(function () {
				seed = null; pw = null; pw2 = null; passphrase = null;
			});
		});
	}

	// secretsOf decrypts a wallet's seed, passphrase and next seed with
	// password, and checks the seed gives the wallet's own xpub or first
	// address. A wallet the classic wallet stored without a password holds
	// them in clear.
	function secretsOf(c, w, password) {
		var get = function (encField, plainField) {
			if (w[encField]) return decrypt(w[encField], password);
			return Promise.resolve(w[plainField] || '');
		};
		if (w.encryptedSeed && !password) return Promise.reject(new Error(msg('password')));
		return Promise.all([get('encryptedSeed', 'seed'), get('encryptedSeedPassphrase', 'seedPassphrase'), get('encryptedNextSeed', 'nextSeed')])
			.catch(function () { throw new Error(msg('wrong-password')); })
			.then(function (s) {
				var sec = { seed: s[0], passphrase: s[1], nextSeed: s[2] };
				if (!sec.seed) throw new Error(msg('wrong-password'));
				if (isBip44(w)) {
					if (call(c, 'bip44AccountXPub', [sec.seed, sec.passphrase]) !== w.accounts[0].accountXpubKey) throw new Error(msg('wrong-password'));
				} else if (step(c, initialState(sec.seed, w.walletType)).address !== w.addresses[0].address) {
					throw new Error(msg('wrong-password'));
				}
				return sec;
			});
	}

	// deterministicChain is the wallet's addresses with their secret keys,
	// one more after them, and the seed that comes after that one.
	function deterministicChain(c, w, sec) {
		var out = [], next = initialState(sec.seed, w.walletType);
		for (var i = 0; i <= w.addresses.length; i++) {
			var a = step(c, next);
			out.push({ address: a.address, secret: a.secret });
			next = a.nextSeed;
		}
		return { addrs: out, after: next };
	}

	// freshOf is the next unused change and receive address of w, neither
	// stored yet. A deterministic wallet has one chain, so both are the same.
	function freshOf(c, w, sec) {
		if (isBip44(w)) {
			var a = w.accounts[0];
			return {
				change: call(c, 'bip44XPubAddresses', [a.accountXpubKey, 1, a.changeAddresses.length, 1])[0].address,
				receive: call(c, 'bip44XPubAddresses', [a.accountXpubKey, 0, a.externalAddresses.length, 1])[0].address,
				changeAt: a.changeAddresses.length, receiveAt: a.externalAddresses.length,
			};
		}
		var ch = deterministicChain(c, w, sec), last = ch.addrs[ch.addrs.length - 1];
		return { change: last.address, receive: last.address, after: ch.after, changeAt: w.addresses.length };
	}

	// keysOf maps each address in addrs to its secret key.
	function keysOf(c, w, sec, addrs) {
		var keys = {};
		if (isBip44(w)) {
			var a = w.accounts[0], where = {};
			a.externalAddresses.forEach(function (e, i) { where[e.address] = [0, i]; });
			a.changeAddresses.forEach(function (e, i) { where[e.address] = [1, i]; });
			addrs.forEach(function (addr) {
				var p = where[addr];
				if (!p) throw new Error(msg('mismatch'));
				var d = call(c, 'bip44Addresses', [sec.seed, sec.passphrase, p[0], p[1], 1])[0];
				if (d.address !== addr) throw new Error(msg('mismatch'));
				keys[addr] = d.secret;
			});
			return keys;
		}
		deterministicChain(c, w, sec).addrs.forEach(function (e) { keys[e.address] = e.secret; });
		addrs.forEach(function (addr) { if (!keys[addr]) throw new Error(msg('mismatch')); });
		return keys;
	}

	// droplets reads a coin amount as millionths, or NaN.
	function droplets(s) {
		var m = /^(\d*)(?:\.(\d{0,6}))?$/.exec((s || '').trim());
		if (!m || (m[1] === '' && !m[2])) return NaN;
		return parseInt(m[1] || '0', 10) * 1e6 + parseInt(((m[2] || '') + '000000').slice(0, 6), 10);
	}

	// formDests are the destinations the form asks for, read here rather
	// than taken from the server: bulk lines when there are any, else the rows.
	function formDests(form) {
		var out = [];
		var bulk = form.elements.bulk ? form.elements.bulk.value.trim() : '';
		if (bulk) {
			bulk.split('\n').forEach(function (line) {
				var f = line.split(',').map(function (x) { return x.trim(); });
				if (f.length >= 2 && f[0]) out.push({ to: f[0], coins: droplets(f[1]) });
			});
			return out;
		}
		var to = form.querySelectorAll('#send-rows [name=to]'), coins = form.querySelectorAll('#send-rows [name=coins]');
		for (var i = 0; i < to.length; i++) {
			if (to[i].value.trim()) out.push({ to: to[i].value.trim(), coins: droplets(coins[i].value) });
		}
		return out;
	}

	// checkSign fails unless the outputs pay each destination the form asked
	// for exactly what it asked, and everything else comes back to the
	// wallet. A move has no form destination: it goes to the fresh address.
	function checkSign(s, w, fresh, want) {
		var own = {};
		summary(w).addresses.forEach(function (a) { own[a.address] = true; });
		own[fresh.change] = true;
		if (!want) {
			if (s.dests.length !== 1 || s.dests[0].to !== fresh.receive) throw new Error(msg('mismatch'));
			want = s.dests;
		}
		var asked = {}, paid = {};
		want.forEach(function (d) {
			if (!(d.coins > 0)) throw new Error(msg('mismatch'));
			asked[d.to] = (asked[d.to] || 0) + d.coins;
		});
		s.outputs.forEach(function (o) {
			if (asked[o.address] !== undefined) paid[o.address] = (paid[o.address] || 0) + o.coins;
			else if (!own[o.address]) throw new Error(msg('mismatch'));
		});
		Object.keys(asked).forEach(function (a) {
			var p = paid[a] || 0;
			if (own[a] ? p < asked[a] : p !== asked[a]) throw new Error(msg('mismatch'));
		});
		if (!s.inputs.length) throw new Error(msg('mismatch'));
	}

	// remember stores the fresh addresses the transaction pays, before it is
	// sent. It fails if the stored wallet has no place for one of them.
	function remember(id, fresh, s, password) {
		var paid = {};
		s.outputs.forEach(function (o) { paid[o.address] = true; });
		// add appends addr at index at. It is false for one that is not paid or already there.
		var add = function (chain, addr, at) {
			if (!paid[addr] || chain.some(function (e) { return e.address === addr; })) return false;
			if (chain.length !== at) throw new Error(msg('mismatch'));
			chain.push({ address: addr });
			return true;
		};
		return withList(function (list) {
			var w = list.filter(shown).find(function (x) { return firstAddress(x) === id; });
			if (!w) throw new Error(msg('missing'));
			if (isHw(w)) return add(w.addresses, fresh.change, fresh.changeAt) ? writeStored(list) : null;
			if (isBip44(w)) {
				var a = w.accounts[0];
				var changed = add(a.changeAddresses, fresh.change, fresh.changeAt);
				if (!add(a.externalAddresses, fresh.receive, fresh.receiveAt) && !changed) return null;
				w.addresses = a.externalAddresses.concat(a.changeAddresses);
				return writeStored(list);
			}
			if (!add(w.addresses, fresh.change, fresh.changeAt)) return null;
			if (w.encryptedSeed) {
				return encrypt(fresh.after, password).then(function (e) {
					w.encryptedNextSeed = e;
					return writeStored(list);
				});
			}
			w.nextSeed = fresh.after;
			return writeStored(list);
		});
	}

	// hwSign has the Skywallet sign s, as the Angular wallet's
	// signWithHardwareWallet: inputs by address index, and outputs to the
	// wallet's own addresses marked so the device does not ask about them.
	function hwSign(session, s) {
		var w = session.w, index = {};
		w.addresses.forEach(function (e, i) { index[e.address] = i; });
		index[session.fresh.change] = session.fresh.changeAt;
		var inputs = s.inputs.map(function (i) {
			if (index[i.address] === undefined) throw new Error(msg('mismatch'));
			return { hash: i.hash, index: index[i.address] };
		});
		var outputs = s.outputs.map(function (o) {
			// The daemon reads coins as a decimal amount, not droplets.
			var c = String(o.coins).padStart(7, '0');
			var out = { address: o.address, coins: c.slice(0, -6) + '.' + c.slice(-6), hours: String(o.hours) };
			if (index[o.address] !== undefined) out.address_index = index[o.address];
			return out;
		});
		return window.SkyHW.sign(firstAddress(w), inputs, outputs).then(function (sigs) {
			var ins = s.inputs.map(function (i) { return { hash: i.hash, address: i.address, secret: '' }; });
			return call(session.c, 'prepareTransactionWithSignatures', [JSON.stringify(ins), JSON.stringify(s.outputs), JSON.stringify(sigs)]);
		});
	}

	function sendFlow(id) {
		var form = $('bw-send'), session = null;
		withList(function (list) {
			var w = list.filter(shown).find(function (x) { return firstAddress(x) === id; });
			if (!w) throw new Error(msg('missing'));
			$('bw-data').value = JSON.stringify(summary(w));
			document.body.dispatchEvent(new CustomEvent('bw-ready'));
			if (isHw(w)) $('bw-send-password').closest('label').hidden = true;
			form.addEventListener('submit', function (e) {
				e.preventDefault();
				if (e.submitter) $('bw-unsigned').value = '';
				setError('');
				var pw = $('bw-send-password').value;
				if (isHw(w)) {
					// The change goes to the device's next address, kept once it is paid.
					var at = w.addresses.length;
					Promise.all([cipher(), window.SkyHW.addresses(at, 1)]).then(function (r) {
						var fresh = { change: r[1][0], receive: r[1][0], changeAt: at, receiveAt: at };
						$('bw-fresh-change').value = fresh.change;
						$('bw-fresh-receive').value = fresh.receive;
						session = { c: r[0], w: w, fresh: fresh, pw: '' };
						htmx.trigger(form, 'bw-review');
					}).catch(function (err) { setError(err.message); });
					return;
				}
				cipher().then(function (c) {
					return secretsOf(c, w, pw).then(function (sec) {
						var fresh = freshOf(c, w, sec);
						$('bw-fresh-change').value = fresh.change;
						$('bw-fresh-receive').value = fresh.receive;
						session = { c: c, w: w, sec: sec, fresh: fresh, pw: pw };
						$('bw-send-password').value = '';
						htmx.trigger(form, 'bw-review');
					});
				}).catch(function (err) { setError(err.message); });
			});
		}).catch(function (err) { setError(err && err.message ? err.message : msg('error')); });

		document.addEventListener('click', function (e) {
			if (!e.target.closest || !e.target.closest('#bw-sign-send')) return;
			var box = document.querySelector('[data-bw-sign]');
			if (!box || !session) return;
			var btn = $('bw-sign-send');
			btn.disabled = true;
			try {
				var s = JSON.parse(box.getAttribute('data-bw-sign'));
				checkSign(s, session.w, session.fresh, form.elements.from ? null : formDests(form));
				if (isHw(session.w)) {
					hwSign(session, s).then(function (raw) {
						return remember(id, session.fresh, s, '').then(function () {
							sent = true;
							htmx.ajax('POST', 'bw/send/inject', { target: '#review', swap: 'innerHTML', values: { id: id, rawtx: raw } });
						});
					}).catch(function (err) {
						btn.disabled = false;
						setError(err && err.message ? err.message : msg('error'));
					});
					return;
				}
				var keys = keysOf(session.c, session.w, session.sec, s.inputs.map(function (i) { return i.address; }));
				var inputs = s.inputs.map(function (i) { return { hash: i.hash, secret: keys[i.address], address: i.address }; });
				var raw = call(session.c, 'prepareTransaction', [JSON.stringify(inputs), JSON.stringify(s.outputs)]);
				keys = null;
				inputs = null;
				// The fresh addresses are stored before the node gets the transaction,
				// so the wallet lists them even if the answer never arrives.
				remember(id, session.fresh, s, session.pw).then(function () {
					sent = true;
					htmx.ajax('POST', 'bw/send/inject', { target: '#review', swap: 'innerHTML', values: { id: id, rawtx: raw } });
				}).catch(function (err) {
					btn.disabled = false;
					setError(err && err.message ? err.message : msg('error'));
				});
			} catch (err) {
				btn.disabled = false;
				setError(err.message);
			}
		});

		// Unsigned asks for the same review without signing. The fresh change
		// address it pays is stored before the transaction is shown.
		document.addEventListener('click', function (e) {
			if (!e.target.closest || !e.target.closest('[data-bw-unsigned]')) return;
			$('bw-unsigned').value = '1';
			form.requestSubmit();
		});
		document.addEventListener('htmx:afterSwap', function (e) {
			if (!e.detail.target || e.detail.target.id !== 'review' || !session) return;
			var box = document.querySelector('[data-bw-unsigned-tx]');
			if (!box) return;
			var done = session;
			session = null;
			try {
				var s = JSON.parse(box.getAttribute('data-bw-unsigned-tx'));
				checkSign(s, done.w, done.fresh, formDests(form));
				remember(id, done.fresh, s, done.pw).catch(function (err) {
					box.hidden = true;
					setError(err && err.message ? err.message : msg('error'));
				});
			} catch (err) {
				box.hidden = true;
				setError(err.message);
			}
		});

		// Once the node took the transaction the form is done: a second send
		// from this page would derive the same fresh address again.
		var sent = false;
		document.addEventListener('htmx:afterSwap', function (e) {
			if (!sent || !e.detail.target || e.detail.target.id !== 'review') return;
			sent = false;
			if (!document.querySelector('#review .alert.ok')) return;
			session = null;
			form.hidden = true;
		});
	}

	// nextSeedAfter is the state after the first count addresses of a
	// deterministic or post-quantum wallet of type with seed.
	function nextSeedAfter(c, seed, count, type) {
		var next = initialState(seed, type);
		for (var i = 0; i < count; i++) next = step(c, next).nextSeed;
		return next;
	}

	// sealSecrets stores the wallet's secrets encrypted under password, as
	// the Angular wallet's setWalletPassword does, and drops the clear ones.
	function sealSecrets(w, sec, password) {
		return Promise.all([
			encrypt(sec.seed, password),
			sec.nextSeed ? encrypt(sec.nextSeed, password) : null,
			sec.passphrase ? encrypt(sec.passphrase, password) : null,
		]).then(function (e) {
			w.encryptedSeed = e[0];
			if (e[1]) w.encryptedNextSeed = e[1]; else delete w.encryptedNextSeed;
			if (e[2]) w.encryptedSeedPassphrase = e[2]; else delete w.encryptedSeedPassphrase;
			delete w.seed;
			delete w.nextSeed;
			delete w.seedPassphrase;
		});
	}

	// walletOps runs the option forms of the browser wallet whose first
	// address is id. Each change is written to storage and the page reloads.
	function walletOps(id) {
		function field(form, name) { var e = form.querySelector('[data-f="' + name + '"]'); return e ? e.value : ''; }
		function change(f, done) {
			return withList(function (list) {
				var w = list.filter(shown).find(function (x) { return firstAddress(x) === id; });
				if (!w) throw new Error(msg('missing'));
				return Promise.resolve(f(w, list)).then(function () { return writeStored(list); });
			}).then(function () {
				if (typeof done === 'function') done = done();
				window.location.href = 'bw/wallet?id=' + encodeURIComponent(id) + '&done=' + done;
			});
		}
		var ops = {
			rename: function (form) {
				var label = field(form, 'label').trim();
				if (!label) return Promise.reject(new Error(msg('label')));
				return change(function (w) { w.label = label; }, 'renamed');
			},
			address: function (form) {
				var n = parseInt(field(form, 'num'), 10);
				if (!(n >= 1 && n <= 100)) return Promise.reject(new Error(msg('num')));
				return cipher().then(function (c) {
					return change(function (w) {
						if (isBip44(w)) {
							var a = w.accounts[0];
							call(c, 'bip44XPubAddresses', [a.accountXpubKey, 0, a.externalAddresses.length, n]).forEach(function (x) {
								a.externalAddresses.push({ address: x.address });
							});
							w.addresses = a.externalAddresses.concat(a.changeAddresses);
							return null;
						}
						var pw = field(form, 'password');
						return secretsOf(c, w, pw).then(function (sec) {
							var next = nextSeedAfter(c, sec.seed, w.addresses.length, w.walletType);
							for (var i = 0; i < n; i++) {
								var g = step(c, next);
								w.addresses.push({ address: g.address });
								next = g.nextSeed;
							}
							if (!w.encryptedSeed) {
								w.nextSeed = next;
								return null;
							}
							return encrypt(next, pw).then(function (e) { w.encryptedNextSeed = e; });
						});
					}, 'added&n=' + n);
				});
			},
			scan: function (form) {
				return cipher().then(function (c) {
					var found = 0;
					return change(function (w) {
						// More is the used addresses after the wallet's last ones,
						// derived from start until a batch of 20 has none in use.
						var more = function (derive, start) {
							return scan(function (s, n) { return derive(start + s, n); }, 0);
						};
						if (isBip44(w)) {
							var a = w.accounts[0];
							var xp = function (chain) {
								return function (s, n) { return call(c, 'bip44XPubAddresses', [a.accountXpubKey, chain, s, n]); };
							};
							return Promise.all([more(xp(0), a.externalAddresses.length), more(xp(1), a.changeAddresses.length)]).then(function (r) {
								r[0].forEach(function (x) { a.externalAddresses.push({ address: x.address }); });
								r[1].forEach(function (x) { a.changeAddresses.push({ address: x.address }); });
								w.addresses = a.externalAddresses.concat(a.changeAddresses);
								found = r[0].length + r[1].length;
							});
						}
						var pw = field(form, 'password');
						return secretsOf(c, w, pw).then(function (sec) {
							var chain = [], next = nextSeedAfter(c, sec.seed, w.addresses.length, w.walletType);
							var derive = function (s, n) {
								while (chain.length < s + n) {
									var g = step(c, next);
									chain.push({ address: g.address, nextSeed: g.nextSeed });
									next = g.nextSeed;
								}
								return chain.slice(s, s + n);
							};
							return more(derive, 0).then(function (kept) {
								found = kept.length;
								if (!found) return null;
								kept.forEach(function (x) { w.addresses.push({ address: x.address }); });
								var after = kept[kept.length - 1].nextSeed;
								if (!w.encryptedSeed) { w.nextSeed = after; return null; }
								return encrypt(after, pw).then(function (e) { w.encryptedNextSeed = e; });
							});
						});
					}, function () { return 'scanned&n=' + found; });
				});
			},
			seed: function (form) {
				return withList(function (list) {
					var w = list.filter(ours).find(function (x) { return firstAddress(x) === id; });
					if (!w) throw new Error(msg('missing'));
					return cipher().then(function (c) { return secretsOf(c, w, field(form, 'password')); });
				}).then(function (sec) {
					var s = form.querySelector('[data-f="seed-out"]'), p = form.querySelector('[data-f="passphrase-out"]');
					s.value = sec.seed;
					s.hidden = false;
					p.value = sec.passphrase || '';
					p.hidden = !sec.passphrase;
					form.querySelector('[data-f="password"]').value = '';
				});
			},
			password: function (form) {
				var pw = field(form, 'password'), nw = field(form, 'new');
				if (!nw) return Promise.reject(new Error(msg('new-password')));
				if (nw !== field(form, 'new2')) return Promise.reject(new Error(msg('password-confirm')));
				return cipher().then(function (c) {
					return change(function (w) {
						return secretsOf(c, w, pw).then(function (sec) {
							if (!isBip44(w) && !sec.nextSeed) sec.nextSeed = nextSeedAfter(c, sec.seed, w.addresses.length, w.walletType);
							return sealSecrets(w, sec, nw);
						});
					}, 'reset');
				});
			},
			reset: function (form) {
				var seed = normalize(field(form, 'seed')), passphrase = field(form, 'passphrase'), nw = field(form, 'new');
				if (!nw) return Promise.reject(new Error(msg('new-password')));
				if (nw !== field(form, 'new2')) return Promise.reject(new Error(msg('password-confirm')));
				return cipher().then(function (c) {
					return change(function (w) {
						var sec = { seed: seed, passphrase: isBip44(w) ? passphrase : '' };
						if (isBip44(w)) {
							if (call(c, 'bip44AccountXPub', [seed, sec.passphrase]) !== w.accounts[0].accountXpubKey) throw new Error(msg('wrong-seed'));
						} else {
							if (step(c, initialState(seed, w.walletType)).address !== w.addresses[0].address) throw new Error(msg('wrong-seed'));
							sec.nextSeed = nextSeedAfter(c, seed, w.addresses.length, w.walletType);
						}
						return sealSecrets(w, sec, nw);
					}, 'reset');
				});
			},
			// hw-address asks the Skywallet for the next addresses of its wallet.
			'hw-address': function (form) {
				var n = parseInt(field(form, 'num'), 10);
				if (!(n >= 1 && n <= 100)) return Promise.reject(new Error(msg('num')));
				return withList(function (list) {
					var w = list.filter(isHw).find(function (x) { return firstAddress(x) === id; });
					if (!w) throw new Error(msg('missing'));
					return window.SkyHW.addresses(0, 1).then(function (a) {
						if (a[0] !== id) throw new Error(msg('wrong-device'));
						return window.SkyHW.addresses(w.addresses.length, n);
					});
				}).then(function (got) {
					return change(function (w) { got.forEach(function (a) { w.addresses.push({ address: a }); }); }, 'added&n=' + got.length);
				});
			},
			// hw-confirm shows an address on the Skywallet's screen to compare.
			'hw-confirm': function (form) {
				var i = parseInt(field(form, 'index'), 10);
				return window.SkyHW.addresses(0, 1).then(function (a) {
					if (a[0] !== id) throw new Error(msg('wrong-device'));
					return window.SkyHW.confirm(i);
				}).then(function () {
					form.querySelector('.op-error').textContent = '';
					form.querySelector('[data-f="confirmed"]').hidden = false;
				});
			},
			delete: function () {
				if (!window.confirm(msg('delete'))) return Promise.resolve();
				return withList(function (list) {
					var kept = list.filter(function (x) { return !(shown(x) && firstAddress(x) === id); });
					return writeStored(kept);
				}).then(function () { window.location.href = './'; });
			},
		};
		document.addEventListener('submit', function (e) {
			var form = e.target.closest && e.target.closest('form[data-bw-op]');
			if (!form) return;
			e.preventDefault();
			var out = form.querySelector('.op-error');
			out.textContent = '';
			var btn = form.querySelector('button[type=submit]');
			btn.disabled = true;
			ops[form.getAttribute('data-bw-op')](form).catch(function (err) {
				out.textContent = err && err.message ? err.message : msg('error');
			}).then(function () { btn.disabled = false; });
		});
	}

	// hexBytes and bytesHex convert between hex and bytes.
	function hexBytes(h) {
		var out = new Uint8Array(h.length / 2);
		for (var i = 0; i < out.length; i++) out[i] = parseInt(h.substr(i * 2, 2), 16);
		return out;
	}
	function bytesHex(b) {
		var s = '';
		for (var i = 0; i < b.length; i++) s += ('0' + b[i].toString(16)).slice(-2);
		return s;
	}

	// signInputs signs job's inputs with keys, as the Angular wallet's joint
	// transaction service does: each signature covers SHA256(inner hash ||
	// input hash) and goes in the input's slot of the encoded transaction
	// (src/coin/transactions.go: length u32, type u8, inner hash, then the
	// u32 counted 65 byte signatures).
	function signInputs(job, keys) {
		var extras = window.SkycoinCipherExtras;
		if (!extras || typeof extras.signHash !== 'function') return Promise.reject(new Error(msg('error')));
		var raw = job.raw;
		return job.inputs.reduce(function (p, input) {
			return p.then(function () {
				var data = new Uint8Array(64);
				data.set(hexBytes(job.inner), 0);
				data.set(hexBytes(input.hash), 32);
				return crypto.subtle.digest('SHA-256', data).then(function (d) {
					var sig = extras.signHash(bytesHex(new Uint8Array(d)), keys[input.address]);
					if (!sig || sig.error || !/^[0-9a-f]{130}$/i.test(sig)) throw new Error(msg('mismatch'));
					var start = (4 + 1 + 32 + 4 + input.index * 65) * 2;
					raw = raw.substring(0, start) + sig.toLowerCase() + raw.substring(start + 130);
				});
			});
		}, Promise.resolve()).then(function () { return raw; });
	}

	// txFlow lets a browser wallet review, sign and make joint transactions.
	// The chosen wallet's public part goes in #bw-data; its keys stay here.
	function txFlow() {
		var sel = $('bw-tx-wallet'), wallets = {};
		withList(function (list) {
			list.filter(ours).forEach(function (w) {
				var id = firstAddress(w);
				wallets[id] = w;
				var o = document.createElement('option');
				o.value = id;
				o.textContent = w.label || id;
				sel.appendChild(o);
			});
		}).catch(function (err) { setError(err && err.message ? err.message : msg('error')); });
		sel.addEventListener('change', function () {
			var w = wallets[sel.value];
			$('bw-data').value = w ? JSON.stringify(summary(w)) : '';
			document.body.dispatchEvent(new CustomEvent('bw-wallet'));
		});
		document.addEventListener('click', function (e) {
			if (!e.target.closest || !e.target.closest('#bw-tx-sign')) return;
			var box = document.querySelector('[data-bw-tx-sign]');
			var w = box && wallets[box.getAttribute('data-bw-wallet')];
			var out = $('bw-tx-error');
			if (!w) { out.textContent = msg('missing'); return; }
			var job = JSON.parse(box.getAttribute('data-bw-tx-sign'));
			var btn = $('bw-tx-sign');
			btn.disabled = true;
			out.textContent = '';
			var pw = $('bw-tx-password').value;
			cipher().then(function (c) {
				return secretsOf(c, w, pw).then(function (sec) {
					$('bw-tx-password').value = '';
					return signInputs(job, keysOf(c, w, sec, job.inputs.map(function (i) { return i.address; })));
				});
			}).then(function (raw) {
				htmx.ajax('POST', 'tx/review', { target: '#tx-review', swap: 'innerHTML', values: { rawtx: raw, wallet: $('bw-data').value } });
			}).catch(function (err) {
				btn.disabled = false;
				out.textContent = err && err.message ? err.message : msg('error');
			});
		});
	}

	// appPasswordFlow sets, changes or removes the app password that
	// encrypts the whole stored list, as the Angular wallet's app password
	// page does (lockedWallets holds the list encrypted, wallets holds it
	// when there is none). The wallets inside stay encrypted either way.
	function appPasswordFlow() {
		var locked = !!localStorage.getItem(LOCKED);
		document.querySelectorAll('[data-bw-app]').forEach(function (f) {
			var op = f.getAttribute('data-bw-app');
			f.closest('details').hidden = locked ? op === 'set' : op !== 'set';
		});
		function val(f, n) { var e = f.querySelector('[data-f="' + n + '"]'); return e ? e.value : ''; }
		var ops = {
			set: function (f) {
				var nw = val(f, 'new');
				if (!nw) throw new Error(msg('new-password'));
				if (nw !== val(f, 'new2')) throw new Error(msg('password-confirm'));
				var plain = localStorage.getItem(WALLETS) || '[]';
				return encrypt(plain, nw).then(function (c) {
					localStorage.setItem(LOCKED, c);
					localStorage.removeItem(WALLETS);
				});
			},
			change: function (f) {
				var nw = val(f, 'new');
				if (!nw) throw new Error(msg('new-password'));
				if (nw !== val(f, 'new2')) throw new Error(msg('password-confirm'));
				return decrypt(localStorage.getItem(LOCKED), val(f, 'current')).catch(function () {
					throw new Error(msg('wrong-password'));
				}).then(function (s) { return encrypt(s, nw); }).then(function (c) { localStorage.setItem(LOCKED, c); });
			},
			remove: function (f) {
				return decrypt(localStorage.getItem(LOCKED), val(f, 'current')).catch(function () {
					throw new Error(msg('wrong-password'));
				}).then(function (s) {
					localStorage.setItem(WALLETS, s);
					localStorage.removeItem(LOCKED);
				});
			},
		};
		document.addEventListener('submit', function (e) {
			var f = e.target.closest && e.target.closest('form[data-bw-app]');
			if (!f) return;
			e.preventDefault();
			var out = f.querySelector('.op-error');
			out.textContent = '';
			Promise.resolve().then(function () { return ops[f.getAttribute('data-bw-app')](f); }).then(function () {
				window.location.href = 'bw/app-password?done=' + f.getAttribute('data-bw-app');
			}).catch(function (err) { out.textContent = err && err.message ? err.message : msg('error'); });
		});
	}

	// The unlock form can erase a locked list whose password is lost, as
	// the Angular wallet's lock screen does.
	document.addEventListener('click', function (e) {
		var b = e.target.closest && e.target.closest('[data-bw-erase]');
		if (!b || !window.confirm(b.getAttribute('data-bw-erase'))) return;
		localStorage.removeItem(LOCKED);
		localStorage.removeItem(WALLETS);
		window.location.reload();
	});

	// Auto-sweep of a wallet kept in this browser (bwsweep.templ). The rule
	// is kept in localStorage under SWEEP and holds no secret. The password
	// stays in this page's memory for the time asked. Every pass runs the
	// node's rule code in the wasm cipher (sweepPlan), signs with
	// prepareTransaction as a send does, and hands the result to the node.
	var SWEEP = 'sweepRules';
	function sweepRules() {
		try { return JSON.parse(localStorage.getItem(SWEEP) || '{}') || {}; } catch (_) { return {}; }
	}
	function ruleActive(r) { return !!r && (r.enabled || !!(r.mappings && r.mappings.length)); }
	function apiGet(path) {
		return fetch(root + 'api/' + path, { credentials: 'same-origin' })
			.then(function (r) { if (!r.ok) throw new Error(r.status + ' ' + r.statusText); return r.json(); });
	}
	function keys(o) { return Object.keys(o); }

	// chainOf reads the unspent outputs of addrs, those an unconfirmed
	// transaction spends, and which addresses sent or appear in a block.
	function chainOf(addrs) {
		if (!addrs.length) return Promise.resolve({ outputs: [], outgoing: {}, sent: {}, seen: {} });
		var q = encodeURIComponent(addrs.join(','));
		return Promise.all([apiGet('v1/outputs?addrs=' + q), apiGet('v1/transactions?verbose=1&confirmed=1&addrs=' + q)]).then(function (r) {
			var want = {}, sent = {}, seen = {}, outgoing = {};
			addrs.forEach(function (a) { want[a] = true; });
			(r[1] || []).forEach(function (t) {
				(t.txn.inputs || []).forEach(function (i) { if (want[i.owner]) { sent[i.owner] = true; seen[i.owner] = true; } });
				(t.txn.outputs || []).forEach(function (o) { if (want[o.dst]) seen[o.dst] = true; });
			});
			(r[0].outgoing_outputs || []).forEach(function (o) { outgoing[o.hash] = true; });
			return { outputs: r[0].head_outputs || [], outgoing: outgoing, sent: sent, seen: seen };
		});
	}

	// sweepAddrs are w's addresses with their place, as sweepPlan takes them.
	function sweepAddrs(w) {
		if (!isBip44(w)) return w.addresses.map(function (e) { return { address: e.address }; });
		var a = w.accounts[0], out = [];
		a.externalAddresses.forEach(function (e, i) { out.push({ address: e.address, child: i, bip44: true }); });
		a.changeAddresses.forEach(function (e, i) { out.push({ address: e.address, change: true, child: i, bip44: true }); });
		return out;
	}

	// nextUnused is the first receive address of w after the last one seen,
	// or the next one to derive, which needs the seed unless w is bip44.
	function nextUnused(c, w, sec, seen) {
		var stored = isBip44(w) ? w.accounts[0].externalAddresses : w.addresses, last = -1;
		stored.forEach(function (e, i) { if (seen[e.address]) last = i; });
		if (last + 1 < stored.length) return { address: stored[last + 1].address, index: last + 1 };
		if (isBip44(w)) {
			return { address: call(c, 'bip44XPubAddresses', [w.accounts[0].accountXpubKey, 0, stored.length, 1])[0].address, index: stored.length, fresh: true };
		}
		if (!sec) throw new Error(sweepMsg('dest-full'));
		var ch = deterministicChain(c, w, sec);
		return { address: ch.addrs[ch.addrs.length - 1].address, index: stored.length, fresh: true, after: ch.after };
	}

	// storeFresh keeps a derived destination in its wallet before anything
	// is sent to it, so the wallet lists the coins.
	function storeFresh(wid, f, password) {
		return withList(function (list) {
			var w = list.filter(ours).find(function (x) { return firstAddress(x) === wid; });
			if (!w) throw new Error(msg('missing'));
			if (isBip44(w)) {
				var a = w.accounts[0];
				if (a.externalAddresses.length !== f.index) throw new Error(msg('mismatch'));
				a.externalAddresses.push({ address: f.address });
				w.addresses = a.externalAddresses.concat(a.changeAddresses);
				return writeStored(list);
			}
			if (w.addresses.length !== f.index) throw new Error(msg('mismatch'));
			w.addresses.push({ address: f.address });
			if (w.encryptedSeed) {
				return encrypt(f.after, password).then(function (e) { w.encryptedNextSeed = e; return writeStored(list); });
			}
			w.nextSeed = f.after;
			return writeStored(list);
		});
	}

	var sweepBox = null;
	function sweepMsg(name) { return (sweepBox && sweepBox.getAttribute('data-msg-' + name)) || name; }

	function sweepFlow(id) {
		var box = sweepBox = document.querySelector('[data-bw-sweep]');
		if (!box || box.getAttribute('data-bound')) return;
		box.setAttribute('data-bound', '1');
		var form = box.querySelector('[data-bw-sweep-form]');
		var f = function (root, name) { return root.querySelector('[data-f="' + name + '"]'); };
		var s = function (name) { return box.querySelector('[data-s="' + name + '"]'); };
		var session = null; // { sec, password, expires }
		var state = { pending: [], confirmed: 0, last: '', txids: [], dests: [], error: '', paused: '', res: null };
		var busy = false, lastSeq = -1, retryAt = 0, others = [];

		function rule() { return sweepRules()[id] || null; }
		function wallet() {
			return withList(function (list) {
				others = list.filter(function (x) { return ours(x) && !x.temporary && firstAddress(x) !== id; });
				var w = list.filter(ours).find(function (x) { return firstAddress(x) === id; });
				if (!w) throw new Error(msg('missing'));
				return w;
			});
		}
		function locked(w) { return !!w.encryptedSeed && !session; }

		function render() {
			var r = rule(), on = ruleActive(r), res = state.res || {};
			var badge = s('badge');
			var active = on && !state.paused && !state.error;
			badge.className = 'badge ' + (active ? 'on' : on ? 'warn-text' : 'off');
			badge.textContent = sweepMsg(active ? 'active' : on ? 'paused' : 'off');
			s('paused').hidden = !(on && state.paused);
			s('paused').textContent = state.paused;
			if (s('unlocked')) {
				s('unlocked').textContent = !session ? sweepMsg('no') : session.expires
					? sweepMsg('until').replace('{time}', new Date(session.expires).toLocaleTimeString()) : sweepMsg('until-page');
				box.querySelector('[data-bw-sweep-unlock]').hidden = !!session;
				s('lock').hidden = !session;
			}
			s('exposed').textContent = String(res.exposed || 0);
			s('uncovered').textContent = String(res.uncovered || 0);
			var c = String(res.coins || 0).padStart(7, '0');
			s('waiting').textContent = s('waiting').getAttribute('data-msg').replace('{coins}', c.slice(0, -6) + '.' + c.slice(-6)).replace('{hours}', String(res.hours || 0));
			s('destinations').textContent = state.dests.join('\n');
			s('last').textContent = state.last || sweepMsg('never');
			s('txids').textContent = state.txids.join('\n');
			s('pending').textContent = state.pending.map(function (p) { return p.txid; }).join('\n');
			s('confirmed').textContent = String(state.confirmed);
			s('error').textContent = state.error;
			s('burn').hidden = on;
		}

		// fill puts the stored rule in the form.
		function fill() {
			var r = rule() || { destination: { kind: 'wallet' } };
			f(form, 'enabled').checked = !!r.enabled;
			var d = r.destination || {};
			var to = d.kind === 'address' ? 'address' : d.wallet ? 'wallet' : 'self';
			form.querySelectorAll('[data-f="to"]').forEach(function (e) { e.checked = e.value === to; });
			f(form, 'to_address').value = d.kind === 'address' ? d.address : '';
			f(form, 'include_spent').checked = r.include_spent !== false;
			if (f(form, 'shared')) f(form, 'shared').checked = (r.shared_accounts || []).indexOf(0) >= 0;
			f(form, 'addresses').value = (r.addresses || []).join('\n');
			f(form, 'burn_ok').checked = false;
			var sel = f(form, 'to_wallet'), mapTo = f(form, 'map_to');
			sel.textContent = '';
			mapTo.querySelectorAll('[data-other]').forEach(function (o) { o.remove(); });
			others.forEach(function (o) {
				var label = (o.label || firstAddress(o)) + ' (' + firstAddress(o) + ')';
				sel.appendChild(new Option(label, firstAddress(o), false, d.wallet === firstAddress(o)));
				var opt = new Option(label, 'wallet:' + firstAddress(o));
				opt.setAttribute('data-other', '1');
				mapTo.appendChild(opt);
			});
			maps = (r.mappings || []).slice();
			drawMaps();
			sel.hidden = !others.length;
			form.querySelector('[data-f="to"][value="wallet"]').closest('label').hidden = !others.length;
			if (maps.length || r.include_spent === false || (r.addresses || []).length) form.querySelector('.sweep-advanced').open = true;
		}
		var maps = [];
		function destText(d) {
			if (d.kind === 'address') return d.address;
			if (!d.wallet) return sweepMsg('self');
			var o = others.find(function (x) { return firstAddress(x) === d.wallet; });
			return o && o.label ? o.label + ' (' + d.wallet + ')' : d.wallet;
		}
		function drawMaps() {
			var ul = s('maps');
			ul.textContent = '';
			if (!maps.length) {
				var none = document.createElement('li');
				none.className = 'dim';
				none.textContent = sweepMsg('none');
				ul.appendChild(none);
			}
			maps.forEach(function (m, i) {
				var li = document.createElement('li');
				var from = document.createElement('span'), arrow = document.createElement('span'), to = document.createElement('span');
				from.className = to.className = 'mono';
				from.textContent = m.from.address;
				arrow.className = 'mi';
				arrow.textContent = 'arrow_forward';
				to.textContent = destText(m.to);
				var b = document.createElement('button');
				b.type = 'button';
				b.className = 'btn sm';
				b.setAttribute('data-bw-sweep-remove', String(i));
				b.textContent = sweepMsg('remove');
				li.append(from, arrow, to, b);
				ul.appendChild(li);
			});
		}

		// read is the rule the form describes. Account destinations and
		// sources are not offered: a browser wallet has one account.
		function read() {
			var to = (form.querySelector('[data-f="to"]:checked') || {}).value;
			var r = { wallet: id, enabled: f(form, 'enabled').checked, destination: { kind: 'wallet' } };
			if (to === 'wallet') {
				if (!f(form, 'to_wallet').value) throw new Error(sweepMsg('destination'));
				r.destination = { kind: 'wallet', wallet: f(form, 'to_wallet').value };
			} else if (to === 'address') {
				var a = f(form, 'to_address').value.trim();
				if (!a) throw new Error(sweepMsg('destination'));
				r.destination = { kind: 'address', address: a };
			}
			if (!f(form, 'include_spent').checked) r.include_spent = false;
			if (f(form, 'shared') && f(form, 'shared').checked) r.shared_accounts = [0];
			var addrs = f(form, 'addresses').value.split(/\s+/).filter(Boolean);
			if (addrs.length) r.addresses = addrs;
			if (maps.length) r.mappings = maps.slice();
			return r;
		}

		// save checks r with the cipher against the wallet and the chain,
		// then stores it and runs a pass.
		function save(r) {
			var before = rule();
			if (ruleActive(r) && !ruleActive(before) && !f(form, 'burn_ok').checked) return Promise.reject(new Error(sweepMsg('burn')));
			return Promise.all([cipher(), wallet()]).then(function (got) {
				var c = got[0], w = got[1], addrs = sweepAddrs(w);
				var fixed = [r.destination].concat((r.mappings || []).map(function (m) { return m.to; }))
					.filter(function (d) { return d.kind === 'address'; }).map(function (d) { return d.address; });
				return chainOf(addrs.map(function (a) { return a.address; }).concat(fixed)).then(function (ch) {
					var bad = c.sweepCheck(JSON.stringify({ rule: r, addresses: addrs, sent: keys(ch.sent), accounts: isBip44(w) ? 1 : 0 }));
					if (bad && bad.error) throw new Error(bad.error);
					var all = sweepRules();
					all[id] = r;
					localStorage.setItem(SWEEP, JSON.stringify(all));
				});
			}).then(function () {
				fill();
				done('saved');
				lastSeq = -1;
				return pass();
			});
		}
		function done(name) {
			var d = s('done');
			d.textContent = name ? sweepMsg(name) : '';
			d.hidden = !name;
		}
		function fail(where, err) {
			var e = where.querySelector('.op-error');
			if (e) e.textContent = err && err.message ? err.message : msg('error');
		}

		function plan(c, req) {
			var out = c.sweepPlan(JSON.stringify(req));
			if (out && out.error) throw new Error(out.error);
			return JSON.parse(out);
		}

		// destinationOf resolves target t of r: an address, this wallet's or
		// another browser wallet's next unused address.
		function destinationOf(c, w, sec, r, t, seen) {
			var d = t === -1 ? r.destination : r.mappings[t].to;
			if (d.kind === 'address') return Promise.resolve({ address: d.address });
			if (d.kind !== 'wallet') return Promise.reject(new Error(sweepMsg('account')));
			if (!d.wallet || d.wallet === id) return Promise.resolve(Object.assign(nextUnused(c, w, sec, seen), { wid: id, own: true }));
			var o = others.find(function (x) { return firstAddress(x) === d.wallet; });
			if (!o) return Promise.reject(new Error(msg('missing')));
			var oa = sweepAddrs(o).map(function (a) { return a.address; });
			return chainOf(oa).then(function (ch) { return Object.assign(nextUnused(c, o, null, ch.seen), { wid: d.wallet }); });
		}

		// inject hands raw to the node through bw/sweep/inject, with a fresh CSRF
		// token as ui2.js takes for every write.
		function inject(raw) {
			return fetch(document.body.getAttribute('data-csrf'), { credentials: 'same-origin' })
				.then(function (r) { return r.ok ? r.json() : {}; }).catch(function () { return {}; })
				.then(function (j) {
					var h = { 'Content-Type': 'application/x-www-form-urlencoded' };
					if (j && j.csrf_token) h['X-CSRF-Token'] = j.csrf_token;
					return fetch(root + 'bw/sweep/inject', { method: 'POST', credentials: 'same-origin', headers: h, body: new URLSearchParams({ rawtx: raw }) });
				})
				.then(function (r) { if (!r.ok) throw new Error(r.status + ' ' + r.statusText); return r.text(); })
				.then(function (html) {
					var e = new DOMParser().parseFromString(html, 'text/html').querySelector('[data-txid]');
					if (!e) throw new Error(msg('error'));
					if (e.getAttribute('data-error')) throw new Error(e.getAttribute('data-error'));
					return e.getAttribute('data-txid');
				});
		}

		// pass sweeps once: everything waiting on exposed addresses goes out
		// in this pass, so opening and unlocking the page catches up.
		function pass() {
			if (busy) return Promise.resolve();
			var r = rule();
			if (!ruleActive(r)) { state.paused = ''; state.res = null; render(); return Promise.resolve(); }
			busy = true;
			state.error = '';
			return Promise.all([cipher(), wallet()]).then(function (got) {
				var c = got[0], w = got[1], addrs = sweepAddrs(w);
				return chainOf(addrs.map(function (a) { return a.address; })).then(function (ch) {
					var head = {}, covered = [];
					ch.outputs.forEach(function (o) { head[o.hash] = true; });
					state.pending = state.pending.filter(function (p) {
						if (p.inputs.every(function (h) { return ch.outgoing[h]; })) { covered = covered.concat(p.inputs); return true; }
						if (p.inputs.every(function (h) { return !head[h]; })) state.confirmed++;
						return false;
					});
					var req = {
						rule: r, addresses: addrs, sent: keys(ch.sent), seen: keys(ch.seen), covered: covered, accounts: isBip44(w) ? 1 : 0,
						outputs: ch.outputs.map(function (o) {
							return { hash: o.hash, address: o.address, coins: droplets(o.coins), hours: o.calculated_hours, contested: !!ch.outgoing[o.hash] && covered.indexOf(o.hash) < 0 };
						}),
					};
					var res = state.res = plan(c, req);
					if (res.gap && res.gap['0'] && isBip44(w)) return fillGap(c, w, res.gap['0']).then(function () { busy = false; return pass(); });
					state.paused = '';
					if (!res.need) return finish(c, w, req, res);
					var coins = String(res.coins).padStart(7, '0');
					if (locked(w)) {
						state.paused = sweepMsg('locked').replace('{coins}', coins.slice(0, -6) + '.' + coins.slice(-6));
						return null;
					}
					var secP = session ? Promise.resolve(session.sec) : secretsOf(c, w, '');
					return secP.then(function (sec) {
						var fresh = [];
						return res.need.reduce(function (p, t) {
							return p.then(function () {
								return destinationOf(c, w, sec, r, t, ch.seen).then(function (d) {
									req.destinations = req.destinations || {};
									req.destinations[String(t)] = d.address;
									if (d.fresh) fresh.push(d);
									if (d.own && d.fresh) req.addresses.push({ address: d.address, child: d.index, bip44: isBip44(w) });
								});
							});
						}, Promise.resolve()).then(function () {
							var outside = keys(req.destinations).map(function (k) { return req.destinations[k]; })
								.filter(function (a) { return !addrs.some(function (x) { return x.address === a; }); });
							return chainOf(outside);
						}).then(function (dc) {
							req.sent = req.sent.concat(keys(dc.sent));
							return finish(c, w, req, plan(c, req), sec, fresh);
						});
					});
				});
			}).catch(function (err) {
				state.error = err && err.message ? err.message : String(err);
				retryAt = Date.now() + 10000;
			}).then(function () { busy = false; render(); });
		}

		// finish signs and sends the sweeps of res.
		function finish(c, w, req, res, sec, fresh) {
			state.res = Object.assign(state.res || {}, { exposed: res.exposed, uncovered: res.uncovered, coins: res.coins, hours: res.hours });
			if (res.errors && res.errors.length) state.error = res.errors.join('\n');
			if (res.no_hours && !(res.sweeps && res.sweeps.length)) {
				var coins = String(res.coins).padStart(7, '0');
				state.paused = sweepMsg('no-hours').replace('{coins}', coins.slice(0, -6) + '.' + coins.slice(-6));
			}
			if (!res.sweeps || !res.sweeps.length) return null;
			var used = {};
			res.sweeps.forEach(function (sw) { used[sw.destination] = true; });
			var store = (fresh || []).filter(function (d) { return used[d.address]; });
			var signed = res.sweeps.map(function (sw) {
				var k = keysOf(c, w, sec, sw.inputs.map(function (i) { return i.address; }));
				var ins = sw.inputs.map(function (i) { return { hash: i.hash, address: i.address, secret: k[i.address] }; });
				return { sw: sw, raw: call(c, 'prepareTransaction', [JSON.stringify(ins), JSON.stringify(sw.outputs)]) };
			});
			return store.reduce(function (p, d) {
				return p.then(function () { return storeFresh(d.wid, d, session ? session.password : ''); });
			}, Promise.resolve()).then(function () {
				var errs = [], txids = [];
				return signed.reduce(function (p, x) {
					return p.then(function () {
						return inject(x.raw).then(function (txid) {
							txids.push(txid);
							state.pending.push({ txid: txid, inputs: x.sw.inputs.map(function (i) { return i.hash; }) });
							var sent = x.sw.outputs[0].coins;
							state.res.coins = Math.max(0, state.res.coins - sent);
							state.res.hours = Math.max(0, state.res.hours - x.sw.fee);
						}, function (err) { errs.push(err.message); });
					});
				}, Promise.resolve()).then(function () {
					if (txids.length) {
						state.last = new Date().toLocaleString();
						state.txids = txids;
					}
					state.dests = keys(used);
					if (errs.length) throw new Error(errs.join('\n'));
				});
			});
		}

		// fillGap stores n more receive addresses of a bip44 wallet whose
		// xPub was shared, so a payment past the stored ones is seen.
		function fillGap(c, w, n) {
			var a = w.accounts[0];
			var more = call(c, 'bip44XPubAddresses', [a.accountXpubKey, 0, a.externalAddresses.length, n]);
			return withList(function (list) {
				var x = list.filter(ours).find(function (y) { return firstAddress(y) === id; });
				var xa = x.accounts[0];
				if (xa.externalAddresses.length !== a.externalAddresses.length) return null;
				more.forEach(function (m) { xa.externalAddresses.push({ address: m.address }); });
				x.addresses = xa.externalAddresses.concat(xa.changeAddresses);
				return writeStored(list);
			});
		}

		// tick runs a pass on every new block, or again after a failure.
		function tick() {
			if (session && session.expires && Date.now() >= session.expires) { session = null; render(); }
			if (!ruleActive(rule()) || busy) return;
			apiGet('v1/blockchain/metadata').then(function (m) {
				var seq = m.head.seq;
				if (seq !== lastSeq || (state.error && Date.now() >= retryAt)) {
					lastSeq = seq;
					return pass();
				}
				return null;
			}).catch(function () { /* the node is away, try on the next tick */ });
		}

		form.addEventListener('submit', function (e) {
			e.preventDefault();
			done('');
			form.querySelector('.op-error').textContent = '';
			Promise.resolve().then(function () { return save(read()); }).catch(function (err) { fail(form, err); });
		});
		box.addEventListener('click', function (e) {
			var add = e.target.closest('[data-bw-sweep-add]'), rm = e.target.closest('[data-bw-sweep-remove]');
			if (!add && !rm && !e.target.closest('[data-bw-sweep-lock]')) return;
			form.querySelector('.op-error').textContent = '';
			done('');
			if (e.target.closest('[data-bw-sweep-lock]')) {
				session = null;
				done('locked-done');
				render();
				return;
			}
			var keep = maps.slice();
			try {
				if (add) {
					var from = f(form, 'map_from').value.trim(), to = f(form, 'map_to').value;
					var dest = to === 'address' ? { kind: 'address', address: f(form, 'map_to_address').value.trim() }
						: to === 'wallet' ? { kind: 'wallet' } : { kind: 'wallet', wallet: to.slice(7) };
					if (!from || (dest.kind === 'address' && !dest.address)) throw new Error(sweepMsg('mapping'));
					maps.push({ from: { address: from }, to: dest });
				} else {
					maps.splice(parseInt(rm.getAttribute('data-bw-sweep-remove'), 10), 1);
				}
				save(read()).then(function () {
					f(form, 'map_from').value = '';
					f(form, 'map_to_address').value = '';
				}, function (err) { maps = keep; drawMaps(); fail(form, err); });
			} catch (err) {
				maps = keep;
				fail(form, err);
			}
		});
		var unlockForm = box.querySelector('[data-bw-sweep-unlock]');
		if (unlockForm) {
			unlockForm.addEventListener('submit', function (e) {
				e.preventDefault();
				unlockForm.querySelector('.op-error').textContent = '';
				var pw = f(unlockForm, 'password').value, ttl = parseInt(f(unlockForm, 'ttl').value, 10) || 0;
				Promise.all([cipher(), wallet()]).then(function (got) {
					return secretsOf(got[0], got[1], pw);
				}).then(function (sec) {
					f(unlockForm, 'password').value = '';
					session = { sec: sec, password: pw, expires: ttl ? Date.now() + ttl * 1000 : 0 };
					done('unlocked');
					lastSeq = -1;
					return pass();
				}).catch(function (err) { fail(unlockForm, err); });
			});
		}
		// The page forgets the password when it is left.
		window.addEventListener('pagehide', function () { session = null; });

		wallet().then(function () { fill(); render(); return pass(); }).catch(function (err) { state.error = err.message; render(); });
		setInterval(tick, 5000);
	}

	// static/hw.js keeps a Skywallet's wallet in the same list.
	window.SkyBW = { withList: withList, writeStored: writeStored, firstAddress: firstAddress, isHw: isHw };

	document.addEventListener('DOMContentLoaded', function () {
		var fail = function (err) {
			var e = $('bw-page-error') || $('bw-error');
			if (e) e.textContent = err && err.message ? err.message : msg('error');
		};
		if (mode === 'list') {
			withList(function (list) { ready(list.filter(shown).map(summary)); }).catch(fail);
		} else if (mode === 'hw') {
			var hwList = function () { withList(function (l) { ready(l.filter(isHw).map(summary)); }).catch(fail); };
			document.body.addEventListener('hw-changed', hwList);
			hwList();
		} else if (mode === 'wallet') {
			var id = page.getAttribute('data-bw-id');
			withList(function (list) {
				var w = list.filter(shown).find(function (x) { return firstAddress(x) === id; });
				if (!w) throw new Error(msg('missing'));
				ready(summary(w));
			}).catch(fail);
			walletOps(id);
			document.body.addEventListener('htmx:afterSwap', function (e) {
				if (e.detail.target && e.detail.target.id === 'bw-wallet') sweepFlow(id);
			});
		} else if (mode === 'new' || mode === 'load') {
			createFlow(mode === 'load');
		} else if (mode === 'send') {
			sendFlow(page.getAttribute('data-bw-id'));
		} else if (mode === 'tx') {
			txFlow();
		} else if (mode === 'app-password') {
			appPasswordFlow();
		}
	});
})();
