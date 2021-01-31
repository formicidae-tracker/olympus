import { HttpClientTestingModule } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { ZoneService } from './zone';

describe('ZoneService', () => {
	beforeEach(() => TestBed.configureTestingModule({
		imports: [ HttpClientTestingModule ],
		providers: [ZoneService]

	}));

  it('should be created', () => {
    const service: ZoneService = TestBed.inject(ZoneService);
    expect(service).toBeTruthy();
  });
});
