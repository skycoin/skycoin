import { CommonTextPipe } from './common-text.pipe';

describe('CommonTextPipe', () => {
  it('create an instance', () => {
    const pipe = new CommonTextPipe(null);
    expect(pipe).toBeTruthy();
  });
});
