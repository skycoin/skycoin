// Loads two browser wallets, creates a joint transaction as the customer, and
// signs it in the browser as each party. ci-scripts/web-wallet-joint-e2e.sh sets
// the JOINT_* variables from the chain it funded and checks the result on it.
import * as fs from 'fs';

const env = (name: string): string => {
  const value = process.env[name];
  if (!value) {
    throw new Error(name + ' is not set');
  }

  return value;
};

const shots = process.env.JOINT_SHOTS || '';
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

describe('web wallet joint transaction', () => {
  it('loads the customer and vendor seeds as browser wallets', async () => {
    await browser.url('/');
    // The first visit asks for a language and acceptance of the disclaimer.
    const english = await $('.//button[contains(normalize-space(.), "English")]');
    await english.waitForClickable();
    await english.click();
    const accept = await $('mat-dialog-container mat-checkbox label');
    await accept.waitForClickable();
    await shot('disclaimer');
    await accept.click();
    await clickButton('Continue', 'mat-dialog-container');
    await waitForDialogsToClose();

    for (const [label, seed] of [['customer', env('JOINT_CUSTOMER_SEED')], ['vendor', env('JOINT_VENDOR_SEED')]]) {
      if (label === 'vendor') {
        await browser.url('/#/wizard/create');
      }
      await clickButton('Load', '.-toggle-container');
      await setText('#label', label);
      await (await $('#wallet_type')).selectByAttribute('value', 'deterministic');
      await setText('#seed', seed);
      await shot('load-' + label);
      await clickButton('Load', '.-buttons-footer');
      await waitFor(async () => (await browser.getUrl()).includes('/wallets'), 'the wallet did not load', 120000);
      await waitForDialogsToClose();
    }

    await shot('wallets');
    // Each browser wallet holds the output funded on the node's wallet of the same seed.
    const page = await $('body').getText();
    expect(page).toMatch(/customer\s+0\.002/);
    expect(page).toMatch(/vendor\s+0\.001/);
  });

  it('creates the transaction as the customer and signs it as each party', async () => {
    await browser.url('/#/joint');
    await (await $('#joint-wallet')).waitForDisplayed();
    await (await $('#joint-wallet')).selectByVisibleText('customer');

    await (await $('#joint-outputs')).click();
    const carrier = await $(`.//mat-option[contains(normalize-space(.), "${env('JOINT_CUSTOMER_UX')}")]`);
    await carrier.waitForClickable();
    await carrier.click();
    await browser.keys('Escape');

    await setText('#joint-theirs', env('JOINT_VENDOR_UX'));
    await (await $('#joint-change')).click();
    await waitFor(async () => (await $('body').getText()).includes(env('JOINT_VENDOR_ADDR')),
      'the vendor output was not looked up');

    const outputs: [string, string, string][] = [
      [env('JOINT_CUSTOMER_BACK'), '0.002', '0'],
      [env('JOINT_VENDOR_PAID'), '0.001', env('JOINT_PAY')],
    ];
    const rows = await $$('.-destination');
    for (let i = 0; i < outputs.length; i++) {
      const inputs = await rows[i].$$('input');
      for (let j = 0; j < 3; j++) {
        await inputs[j].setValue(outputs[i][j]);
      }
    }
    await shot('create');
    await clickButton('Create', '.-paper .-buttons');
    await waitFor(async () => (await $('#joint-created')).isDisplayed(), 'the transaction was not created');
    expect(await (await $('#joint-created')).getValue()).toMatch(/^[0-9a-f]+$/);

    await clickButton('Review', '.-paper .-buttons');
    await waitFor(async () => (await $('body').getText()).includes('0 of 2 inputs signed'), 'the review did not show');
    await shot('review');

    await clickButton('Sign my inputs', '.-paper .-buttons');
    await waitFor(async () => (await $('body').getText()).includes('1 of 2 inputs signed'), 'the customer did not sign');
    await shot('customer-signed');
    expect(await $('body').getText()).toContain('This wallet has already signed its inputs.');

    await (await $('#joint-sign-wallet')).selectByVisibleText('vendor');
    const sign = await $('.-paper .-buttons').$('.//app-button[contains(normalize-space(.), "Sign my inputs")]//button');
    await waitFor(async () => sign.isEnabled(), 'the vendor could not sign');
    await sign.click();
    await waitFor(async () => (await $('body').getText()).includes('2 of 2 inputs signed'), 'the vendor did not sign');
    await shot('vendor-signed');

    const signed = await (await $('#joint-raw')).getValue();
    fs.writeFileSync(env('JOINT_OUT_FILE'), signed.trim());

    await clickButton('Broadcast', '.-paper .-buttons');
    await waitFor(async () => (await (await $('#joint-raw')).getValue()) === '', 'the broadcast did not finish');
    await shot('broadcast');
  });
});
