import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { VideoJsComponent } from './video-js.component';

describe('VideoJsComponent', () => {
	let component: VideoJsComponent;
	let fixture: ComponentFixture<VideoJsComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [VideoJsComponent]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(VideoJsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
