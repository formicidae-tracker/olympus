import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { OlympusService } from './olympus.service';
import { ZoneReport } from '../zone-report';

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
        expect(reports.length).toBe(fakeDB._api_zones.length);
      });
      const req = httpMock.expectOne('/api/zones');
      expect(req.request.method).toBe('GET');
      req.flush(fakeDB._api_zones);
    });
  });

  describe('getZoneReport', () => {
    it('should call the right endpoint', () => {
      service.getZoneReport('minerva', 'box').subscribe((report) => {
        expect(report).toEqual(
          ZoneReport.fromPlain(fakeDB._api_host_minerva_zone_box)
        );
      });
      const req = httpMock.expectOne('/api/host/minerva/zone/box');
      expect(req.request.method).toBe('GET');
      req.flush(fakeDB._api_host_minerva_zone_box);
    });
  });
});
