import { HttpClientTestingModule } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { OlympusService } from './olympus';

describe('OlympusService', () => {
	beforeEach(() => TestBed.configureTestingModule({
		imports: [ HttpClientTestingModule ],
		providers: [OlympusService]

	}));

  it('should be created', () => {
    const service: OlympusService = TestBed.inject(OlympusService);
    expect(service).toBeTruthy();
  });
});
