// Sends from a deterministic browser wallet twice: signed in the browser, then
// as an unsigned transaction signed in the offline dialog. Each time the wallet
// must list the address the change went to and show all its coins.
// ci-scripts/web-wallet-send-e2e.sh sets the SEND_* variables from its chain.

const env = (name: string): string => {
  const value = process.env[name];
  if (!value) {
    throw new Error(name + ' is not set');
  }

  return value;
};

const shots = process.env.SEND_SHOTS || '';
async function shot(name: string) {
  if (shots) {
    await browser.saveScreenshot(`${shots}/${name}.png`);
  }
}

/** Clicks the app-button whose text contains text, inside scope. */
async function clickButton(text: string, scope = 'body') {
  const button = await $(scope).$(`.//app-button[contains(normalize-space(.), "${text}")]//button`);
  await button.waitForExist();
  await button.scrollIntoView({ block: 'center' });
  await button.waitForClickable();
  await button.click();
}

/** Waits for condition, saving a screenshot when it never holds. */
async function waitFor(condition: () => Promise<boolean>, message: string, timeout = 30000) {
  try {
    await browser.waitUntil(condition, { timeout, timeoutMsg: message });
  } catch (e) {
    await shot('failed');
    console.log(await $('body').getText());
    throw e;
  }
}

async function waitForDialogsToClose() {
  await browser.waitUntil(async () => (await $$('mat-dialog-container').length) === 0, { timeoutMsg: 'a dialog stayed open' });
}

async function setText(selector: string, value: string) {
  const field = await $(selector);
  await field.waitForDisplayed();
  await field.setValue(value);
}

/** The addresses the stored wallet lists. */
async function storedAddresses(): Promise<string[]> {
  return browser.execute(() => JSON.parse(localStorage.getItem('wallets') || '[]')[0].addresses.map((a: any) => a.address));
}

/** Selects the wallet in a send form once its balance has loaded. */
async function selectWallet(selector: string) {
  const option = await $(selector).$('.//option[contains(normalize-space(.), "sender")]');
  await waitFor(async () => (await option.isExisting()) && (await option.isEnabled()), 'the wallet has no balance to send', 60000);
  await (await $(selector)).selectByIndex(1);
}

/** The coins the wallets page shows for the wallet. */
async function shownCoins(): Promise<number> {
  const match = (await $('body').getText()).match(/sender\s+([\d.]+)/);

  return match ? Number(match[1]) : NaN;
}

/**
 * Waits for the send to confirm, then checks that the wallet lists exactly
 * listed, whose last is where the change went, and shows coins.
 */
async function expectWallet(coins: number, listed: string[]) {
  await browser.url('/#/wallets');
  const change = listed[listed.length - 1];
  await waitFor(async () => {
    const outputs: any = await (await fetch(`${env('SEND_NODE')}/api/v1/outputs?addrs=${change}`)).json();

    return outputs.head_outputs.length > 0;
  }, 'the transaction did not confirm', 120000);

  await waitFor(async () => (await shownCoins()) === coins, `the wallet does not show ${coins} coins`, 90000);
  expect(await storedAddresses()).toEqual(listed);

  // The row stays open from the last look, and a click would close it.
  if (!(await $('body').getText()).includes(listed[0])) {
    await (await $('.-wallet')).click();
  }
  await waitFor(async () => (await $('body').getText()).includes(change), 'the page does not list the change address');
  await shot('wallet-' + coins);
}

describe('web wallet send', () => {
  const addresses = env('SEND_ADDRS').split(',');

  it('loads the seed as a deterministic browser wallet with one address', async () => {
    await browser.url('/');
    // The first visit asks for a language and acceptance of the disclaimer.
    const english = await $('.//button[contains(normalize-space(.), "English")]');
    await english.waitForClickable();
    await english.click();
    const accept = await $('mat-dialog-container mat-checkbox label');
    await accept.waitForClickable();
    await accept.click();
    await clickButton('Continue', 'mat-dialog-container');
    await waitForDialogsToClose();

    await clickButton('Load', '.-toggle-container');
    await setText('#label', 'sender');
    await (await $('#wallet_type')).selectByAttribute('value', 'deterministic');
    await setText('#seed', env('SEND_SEED'));
    await clickButton('Load', '.-buttons-footer');
    await waitFor(async () => (await browser.getUrl()).includes('/wallets'), 'the wallet did not load', 120000);
    await waitForDialogsToClose();

    await expectWallet(10, addresses.slice(0, 1));
  });

  it('lists the change address of a send signed in the browser', async () => {
    await browser.url('/#/send');
    await selectWallet('#wallet');
    await setText('#address', env('SEND_TO'));
    await setText('#amount', '3');
    await clickButton('Verify', '.-paper .-buttons');
    await shot('verify');
    await clickButton('Send', '.-paper .-buttons');
    // The form comes back once the node has answered.
    await waitFor(async () => (await $('#amount')).isDisplayed(), 'the send did not finish');

    await expectWallet(7, addresses.slice(0, 2));
  });

  it('lists the change address of an unsigned transaction signed in the offline dialog', async () => {
    const dialog = 'mat-dialog-container';
    const rawTx = async () => (await $('#offline-raw-tx')).getValue();
    const hex = /^[0-9a-f]+$/;

    await browser.url('/#/send');
    await clickButton('Advanced');
    await selectWallet('#wallets');
    await setText('#destination0', env('SEND_TO'));
    await setText('#amount0', '2');
    await clickButton('Unsigned', '.-paper .-buttons');
    await waitFor(async () => hex.test(await rawTx()), 'the unsigned transaction was not shown');
    const unsigned = await rawTx();
    await shot('unsigned');
    await clickButton('Close', dialog);
    await waitForDialogsToClose();

    await (await $('.//span[contains(normalize-space(.), "Sign a transaction")]')).click();
    await setText('#offline-raw-tx', unsigned);
    await clickButton('Sign', dialog);
    await waitFor(async () => hex.test(await rawTx()) && (await rawTx()) !== unsigned, 'the transaction was not signed');
    const signed = await rawTx();
    await clickButton('Close', dialog);
    await waitForDialogsToClose();

    await (await $('.//span[contains(normalize-space(.), "Broadcast a transaction")]')).click();
    await setText('#offline-raw-tx', signed);
    await clickButton('Broadcast', dialog);
    await waitForDialogsToClose();

    await expectWallet(5, addresses);
  });
});
