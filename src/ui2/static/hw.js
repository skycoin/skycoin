// The Skywallet hardware wallet, reached through the Skywallet daemon on the
// machine running this browser, as the Angular wallet does it
// (src/skycoin-web/src/app/services/hw-wallet-daemon.service.ts). The device
// keeps the seed and signs; nothing here sees a key.
(function () {
	'use strict';

	var DAEMON = 'http://127.0.0.1:9510/api/v1/';
	var dlg = null;

	function $(id) { return document.getElementById(id); }
	function text(name) {
		var t = document.querySelector('[data-hw-texts]');
		return (t && t.getAttribute('data-msg-' + name)) || name;
	}

	// errorText maps a daemon message to one the page can show, as
	// getHardwareWalletErrorType does.
	function errorText(m) {
		var u = (m || '').toUpperCase();
		var map = [
			[['FAILED OR REFUSED', 'CANCELED BY USER', 'CANCELLED BY USER', 'EXPECTED WORDACK AFTER BUTTON'], 'refused'],
			[['PIN INVALID'], 'wrong-pin'],
			[['WRONG WORD RETYPED'], 'wrong-word'],
			[['PIN MISMATCH'], 'pin-mismatch'],
			[['MNEMONIC NOT SET', 'MNEMONIC REQUIRED'], 'no-seed'],
			[['INVALID SEED, ARE WORDS IN CORRECT ORDER?'], 'invalid-seed'],
			[['THE SEED IS VALID BUT DOES NOT MATCH'], 'wrong-seed'],
			[['LIBUSB', 'FAILED TO FETCH', 'NETWORKERROR'], 'daemon'],
			[['HIDAPI', 'DEVICE DISCONNECTED', 'NO DEVICE CONNECTED'], 'disconnected'],
		];
		for (var i = 0; i < map.length; i++) {
			if (map[i][0].some(function (s) { return u.indexOf(s) >= 0; })) return text('err-' + map[i][1]);
		}
		return m || text('err-generic');
	}
	function fail(m) { var e = new Error(errorText(m)); e.raw = m; throw e; }

	function request(method, route, body) {
		var opts = { method: method, headers: {} };
		if (body !== undefined) {
			opts.headers['Content-Type'] = 'application/json';
			opts.body = JSON.stringify(body);
		}
		return fetch(DAEMON + route, opts).catch(function (err) { fail('NetworkError ' + err.message); }).then(function (r) {
			return r.json().catch(function () { return {}; }).then(function (j) {
				if (!r.ok || j.error) fail((j.error && j.error.message) || r.status + ' ' + r.statusText);
				var d = j.data;
				if (Array.isArray(d) && d.length === 1) d = d[0];
				return d;
			});
		});
	}

	// The dialog asks for what the device wants mid operation: the PIN on
	// its scrambled matrix, a seed word, the passphrase, or a button press.
	function dialog() {
		if (dlg) return dlg;
		dlg = $('hw-dialog');
		dlg.addEventListener('click', function (e) {
			var k = e.target.closest && e.target.closest('[data-hw-key]');
			if (k) {
				var pin = $('hw-pin');
				if (pin.value.length < 9) pin.value += k.getAttribute('data-hw-key');
			}
			if (e.target.closest && e.target.closest('[data-hw-del]')) $('hw-pin').value = $('hw-pin').value.slice(0, -1);
		});
		return dlg;
	}
	// ask shows the part named kind and resolves with what was typed, or
	// null when cancelled.
	function ask(kind, title) {
		var d = dialog();
		['pin', 'word', 'passphrase', 'button'].forEach(function (k) { $('hw-ask-' + k).hidden = k !== kind; });
		$('hw-dialog-title').textContent = title;
		$('hw-dialog-ok').hidden = kind === 'button';
		$('hw-pin').value = '';
		$('hw-word').value = '';
		$('hw-passphrase').value = '';
		if (!d.open) d.showModal();
		if (kind === 'button') return Promise.resolve('');
		var input = { pin: null, word: $('hw-word'), passphrase: $('hw-passphrase') }[kind];
		if (input) input.focus();
		return new Promise(function (resolve) {
			d.onclose = function () { resolve(d.returnValue === 'ok' ? (kind === 'pin' ? $('hw-pin').value : input.value) : null); };
			d.returnValue = '';
		});
	}
	function closeDialog() { if (dlg && dlg.open) { dlg.onclose = null; dlg.close(); } }

	var pinTitle = 'pin';
	// run calls the daemon and answers its requests until the operation ends.
	function run(method, route, body) {
		function step(d) {
			if (typeof d !== 'string') return d;
			var next = null;
			if (d.indexOf('PinMatrixRequest') >= 0) {
				next = ask('pin', text(pinTitle)).then(function (pin) {
					if (pinTitle === 'pin-new') pinTitle = 'pin-confirm';
					return pin ? request('POST', 'intermediate/pin_matrix', { pin: pin }) : cancel();
				});
			} else if (d.indexOf('WordRequest') >= 0) {
				next = ask('word', text('word')).then(function (w) {
					return w ? request('POST', 'intermediate/word', { word: w.trim().toLowerCase() }) : cancel();
				});
			} else if (/passphraserequest/i.test(d)) {
				next = ask('passphrase', text('passphrase')).then(function (p) {
					return p !== null ? request('POST', 'intermediate/passphrase', { passphrase: p }) : cancel();
				});
			} else if (d.indexOf('ButtonRequest') >= 0) {
				next = ask('button', text('button')).then(function () { return request('POST', 'intermediate/button'); });
			}
			return next ? next.then(step) : d;
		}
		return request(method, route, body).then(step).then(function (d) {
			closeDialog();
			return d;
		}, function (err) {
			closeDialog();
			throw err;
		});
	}
	function cancel() {
		return request('PUT', 'cancel').then(function () { fail('Cancelled by user'); });
	}

	// connected is the device's features, or null when none is connected.
	// The daemon answers /available on USB only, so the features decide.
	function features() {
		return request('GET', 'features').catch(function (err) {
			if (/no device connected|device disconnected|hidapi|libusb|emulator/i.test(err.raw || '')) return null;
			throw err;
		});
	}

	var HW = {
		features: features,
		// addresses are count device addresses from index start.
		addresses: function (start, count) {
			return request('PUT', 'cancel').catch(function () {}).then(function () {
				return run('POST', 'generate_addresses', { address_n: count, start_index: start, confirm_address: false });
			}).then(function (d) { return Array.isArray(d) ? d : [d]; });
		},
		confirm: function (index) {
			return run('POST', 'generate_addresses', { address_n: 1, start_index: index, confirm_address: true });
		},
		// sign checks the device is the wallet's, then has it sign. inputs are
		// {hash, index} and outputs {address, coins, hours, address_index?}.
		sign: function (first, inputs, outputs) {
			return HW.addresses(0, 1).then(function (a) {
				if (a[0] !== first) fail(text('err-wrong-device'));
				pinTitle = 'pin';
				return run('POST', 'transaction_sign', { transaction_inputs: inputs, transaction_outputs: outputs });
			}).then(function (d) { return Array.isArray(d) ? d : [d]; });
		},
		errorText: errorText,
	};
	window.SkyHW = HW;

	// The Skywallet page: what the device reports and what it can do.
	function devicePage(page) {
		var info = $('hw-info'), out = $('hw-error'), done = $('hw-done');
		function show(state) {
			page.querySelectorAll('[data-hw-state]').forEach(function (e) {
				e.hidden = e.getAttribute('data-hw-state').split(' ').indexOf(state) < 0;
			});
		}
		var dev = null;
		function refresh() {
			show('checking');
			return features().then(function (f) {
				dev = f;
				if (!f) return show('disconnected');
				$('hw-label').textContent = f.label || '';
				$('hw-label-input').value = f.label || '';
				$('hw-firmware').textContent = [f.fw_major, f.fw_minor, f.fw_patch].join('.');
				$('hw-pin-set').textContent = text(f.pin_protection ? 'yes' : 'no');
				$('hw-passphrase-set').textContent = text(f.passphrase_protection ? 'yes' : 'no');
				$('hw-backup-warn').hidden = !f.needs_backup;
				$('hw-pin-warn').hidden = !f.initialized || !!f.pin_protection;
				page.querySelectorAll('[data-hw-pin]').forEach(function (b) {
					b.hidden = b.getAttribute('data-hw-pin') !== (f.pin_protection ? 'set' : 'unset');
				});
				$('hw-passphrase-toggle').textContent = text(f.passphrase_protection ? 'passphrase-off' : 'passphrase-on');
				show(f.bootloader_mode ? 'bootloader' : f.initialized ? 'ready' : 'empty');
			}).catch(function (err) {
				out.textContent = err.message;
				show('disconnected');
			});
		}
		function op(f, ok) {
			out.textContent = '';
			done.textContent = '';
			show('working');
			return f().then(function () {
				done.textContent = text(ok);
			}, function (err) {
				out.textContent = err.message;
			}).then(refresh);
		}
		var ops = {
			retry: refresh,
			'generate-12': function () { return op(function () { return run('POST', 'generate_mnemonic', { word_count: 12, use_passphrase: false }); }, 'seed-generated'); },
			'generate-24': function () { return op(function () { return run('POST', 'generate_mnemonic', { word_count: 24, use_passphrase: false }); }, 'seed-generated'); },
			'restore-12': function () { return op(function () { return run('POST', 'recovery', { word_count: 12, dry_run: false }); }, 'seed-restored'); },
			'restore-24': function () { return op(function () { return run('POST', 'recovery', { word_count: 24, dry_run: false }); }, 'seed-restored'); },
			'check-12': function () { return op(function () { return run('POST', 'recovery', { word_count: 12, dry_run: true }); }, 'seed-matches'); },
			'check-24': function () { return op(function () { return run('POST', 'recovery', { word_count: 24, dry_run: true }); }, 'seed-matches'); },
			backup: function () { return op(function () { return run('POST', 'backup'); }, 'backed-up'); },
			'set-pin': function () {
				pinTitle = dev && dev.pin_protection ? 'pin-current' : 'pin-new';
				return op(function () {
					return run('POST', 'configure_pin_code', {}).then(function (d) { pinTitle = 'pin'; return d; }, function (e) { pinTitle = 'pin'; throw e; });
				}, 'pin-changed');
			},
			'remove-pin': function () { return op(function () { return run('POST', 'configure_pin_code', { remove_pin: true }); }, 'pin-removed'); },
			passphrase: function () {
				var on = !(dev && dev.passphrase_protection);
				return op(function () { return run('POST', 'apply_settings', { use_passphrase: on }); }, 'settings-applied');
			},
			wipe: function () {
				if (!window.confirm(text('wipe-confirm'))) return Promise.resolve();
				return op(function () { return run('DELETE', 'wipe'); }, 'wiped');
			},
			add: function () { return op(addWallet, 'added'); },
		};
		page.addEventListener('click', function (e) {
			var b = e.target.closest && e.target.closest('[data-hw-op]');
			if (!b) return;
			e.preventDefault();
			ops[b.getAttribute('data-hw-op')]();
		});
		$('hw-label-form').addEventListener('submit', function (e) {
			e.preventDefault();
			var label = $('hw-label-input').value.trim();
			if (!label || (dev && label === dev.label)) return;
			op(function () { return run('POST', 'apply_settings', { label: label }); }, 'label-changed');
		});
		refresh();
	}

	// addWallet keeps the device's wallet in this browser, in the Angular
	// wallet's format, so both UIs list it.
	function addWallet() {
		var BW = window.SkyBW;
		return HW.addresses(0, 1).then(function (a) {
			return features().then(function (f) {
				return BW.withList(function (list) {
					if (list.some(function (w) { return BW.isHw(w) && BW.firstAddress(w) === a[0]; })) fail(text('err-already-added'));
					list.push({ label: (f && f.label) || text('default-label'), addresses: [{ address: a[0] }], isHardware: true, coinId: 0 });
					return BW.writeStored(list);
				});
			});
		}).then(function () { document.body.dispatchEvent(new CustomEvent('hw-changed')); });
	}

	document.addEventListener('DOMContentLoaded', function () {
		var page = document.querySelector('[data-hw-page]');
		if (page) devicePage(page);
	});
})();
