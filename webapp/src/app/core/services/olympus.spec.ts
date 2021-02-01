import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { environment } from '@environments/environment';

import { OlympusService } from './olympus';

describe('OlympusService', () => {
	let httpMock: HttpTestingController;
	let service: OlympusService
	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [OlympusService]
		})
		service = TestBed.inject(OlympusService);
		httpMock = TestBed.inject(HttpTestingController);
	});


	it('should be created', () => {
		expect(service).toBeTruthy();
	});

	it('should fetch zone summaries', () => {
		service.zoneSummaries()
			.subscribe(
				(summaries)=> {
					expect(summaries.length).toBe(3);
				});

		const req = httpMock.expectOne(environment.apiEndpoint+'/zones');
		expect(req.request.method).toBe("GET");
		req.flush('[{"Host":"atreides","Name":"box","Climate":{"Temperature":26.019226705867965,"Humidity":61.4305448316948,"TemperatureBounds":{"Min":20,"Max":31},"HumidityBounds":{"Min":40,"Max":80},"ActiveWarnings":1,"ActiveEmergencies":1,"NumAux":0,"Current":{"Name":"day","Temperature":26,"Humidity":60,"Wind":100,"VisibleLight":40,"UVLight":100},"CurrentEnd":null,"Next":{"Name":"day to night","Temperature":26,"Humidity":60,"Wind":100,"VisibleLight":40,"UVLight":100},"NextEnd":{"Name":"night","Temperature":22,"Humidity":60,"Wind":100,"VisibleLight":0,"UVLight":0},"NextTime":"2021-02-07T17:00:00Z"},"Stream":null},{"Host":"fremens","Name":"box","Climate":{"Temperature":25.992705321693155,"Humidity":60.49926955088248,"TemperatureBounds":{"Min":20,"Max":31},"HumidityBounds":{"Min":40,"Max":80},"ActiveWarnings":1,"ActiveEmergencies":1,"NumAux":0,"Current":{"Name":"day","Temperature":26,"Humidity":60,"Wind":100,"VisibleLight":40,"UVLight":100},"CurrentEnd":null,"Next":{"Name":"day to night","Temperature":26,"Humidity":60,"Wind":100,"VisibleLight":40,"UVLight":100},"NextEnd":{"Name":"night","Temperature":22,"Humidity":60,"Wind":100,"VisibleLight":0,"UVLight":0},"NextTime":"2021-02-07T17:00:00Z"},"Stream":null},{"Host":"fremens","Name":"tunnel","Climate":{"Temperature":-1000,"Humidity":-1000,"TemperatureBounds":{"Min":null,"Max":null},"HumidityBounds":{"Min":null,"Max":null},"ActiveWarnings":0,"ActiveEmergencies":1,"NumAux":0,"Current":{"Name":"always-on","Temperature":-1000,"Humidity":-1000,"Wind":-1000,"VisibleLight":100,"UVLight":0},"CurrentEnd":null,"Next":null,"NextEnd":null,"NextTime":null},"Stream":null}]');

		httpMock.verify();

	});
});
