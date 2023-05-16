import { cases } from 'jasmine-parameterized';
import { humanize_bytes } from './humanize';

describe('humanize_bytes', () => {
  cases([
    [103, '103.0 B'],
    [Math.round(1.029898 * 1024), '1.0 kiB'],
    [Math.round(-13.589898 * 1024 * 1024), '-13.6 MiB'],
    [Math.round(234.12 * 1024 * 1024 * 1024), '234.1 GiB'],
    [Math.round(2.89 * 1024 * 1024 * 1024 * 1024), '2.9 TiB'],
    [Math.round(1024 * 1024 * 1024 * 1024 * 1024), '1.0 PiB'],
  ]).it('should humanize accordingly', ([value, expected]) => {
    expect(humanize_bytes(value)).toEqual(expected);
  });

  it('should use the prefix if given', () => {
    expect(humanize_bytes(1024)).toEqual('1.0 kiB');
    expect(humanize_bytes(1024, 'B/s')).toEqual('1.0 kiB/s');
  });
});
