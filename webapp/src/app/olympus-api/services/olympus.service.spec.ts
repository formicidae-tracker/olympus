import { TestBed } from '@angular/core/testing';

import { OlympusService } from './olympus.service';

describe('OlympusService', () => {
  let service: OlympusService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(OlympusService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
