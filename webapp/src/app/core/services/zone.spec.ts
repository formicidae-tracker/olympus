import { HttpClientTestingModule } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { ZoneClimateReportAdapter } from '@models/zone-climate-report';
import { ZoneSummaryReportAdapter } from '@models/zone-summary-report';
import { ZoneClimateStatusAdapter } from '@models/zone-climate-status'
import { BoundsAdapter } from '@models/bounds';
import { StateAdapter } from '@models/state';
import { ZoneService } from './zone';

describe('ZoneService', () => {
	let statusAdapter = new ZoneClimateStatusAdapter(new BoundsAdapter())
	let climateAdapter = new ZoneClimateReportAdapter(statusAdapter,new StateAdapter());
	let summaryAdapter = new ZoneSummaryReportAdapter(statusAdapter);
	beforeEach(() => TestBed.configureTestingModule({
		imports: [ HttpClientTestingModule ],
		providers: [ZoneService,
					{ provide: ZoneClimateReportAdapter,useValue:climateAdapter},
					{ provide: ZoneSummaryReportAdapter,useValue:summaryAdapter},
				   ]

	}));

  it('should be created', () => {
    const service: ZoneService = TestBed.inject(ZoneService);
    expect(service).toBeTruthy();
  });
});
