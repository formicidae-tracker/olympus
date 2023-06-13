import { cases } from 'jasmine-parameterized';

import { TestBed } from '@angular/core/testing';

import { HumanizeService } from './humanize.service';

describe('HumanizeDurationService', () => {
  let service: HumanizeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(HumanizeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  cases([
    [103, '103.0 B'],
    [Math.round(1.029898 * 1024), '1.0 KiB'],
    [Math.round(-13.589898 * 1024 * 1024), '-13.6 MiB'],
    [Math.round(234.12 * 1024 * 1024 * 1024), '234.1 GiB'],
    [Math.round(2.89 * 1024 * 1024 * 1024 * 1024), '2.9 TiB'],
    [Math.round(1024 * 1024 * 1024 * 1024 * 1024), '1.0 PiB'],
  ]).it('should humanize accordingly', ([value, expected]) => {
    expect(service.humanizeBytes(value)).toEqual(expected);
  });
});
