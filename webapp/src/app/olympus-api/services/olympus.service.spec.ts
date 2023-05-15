import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { OlympusService } from './olympus.service';

import fakeDB from '../fake-backend/db.json';

describe('OlympusService', () => {
  let httpMock: HttpTestingController;
  let service: OlympusService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [OlympusService],
    });
    service = TestBed.inject(OlympusService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getZoneReportSummaries', () => {
    it('should call the right endpoint', () => {
      service.getZoneReportSummaries().subscribe((reports) => {
        expect(reports.length).toBe(3);
      });
      const req = httpMock.expectOne('/api/zones');
      expect(req.request.method).toBe('GET');
      req.flush(fakeDB._api_zones);
    });
  });
});
