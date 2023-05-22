import { TestBed } from '@angular/core/testing';

import { HumanizeDurationService } from './humanize-duration.service.ts~';

describe('HumanizeDurationService', () => {
  let service: HumanizeDurationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(HumanizeDurationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  // No more test, it is simply a wrapper around the library interface
  // to use as a singleton service were needed.
});
