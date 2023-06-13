import { TestBed, ComponentFixture } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { CoreModule } from './core/core.module';
import { NavToolbarComponent } from './nav-toolbar/nav-toolbar.component';
import { ThemeService } from './core/services/theme.service';
import { PushNotificationService } from './core/services/push-notification.service';
import { of } from 'rxjs';

describe('AppComponent', () => {
  let service: ThemeService;
  let push: jasmine.SpyObj<PushNotificationService> =
    jasmine.createSpyObj<PushNotificationService>('PushNotificationService', [
      'updateNotificationsOnDemand',
      'requestSubscriptionOnDemand',
    ]);
  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      imports: [RouterTestingModule, CoreModule],
      declarations: [AppComponent, NavToolbarComponent],
      providers: [{ provide: PushNotificationService, useValue: push }],
    });

    push.updateNotificationsOnDemand.and.callFake(() => of());
    push.requestSubscriptionOnDemand.and.callFake(() => of());

    service = TestBed.inject(ThemeService);
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  describe('render', () => {
    let compiled: HTMLElement;
    let fixture: ComponentFixture<AppComponent>;

    beforeEach(() => {
      fixture = TestBed.createComponent(AppComponent);
      fixture.detectChanges();
      compiled = fixture.nativeElement as HTMLElement;
    });

    it('should apply mode accordingly', () => {
      let classes = compiled
        .querySelector('div')
        ?.attributes.getNamedItem('class')?.textContent;

      expect(classes).toContain('mat-app-background');
      expect(classes).not.toContain('dark-theme');

      service.darkTheme = true;
      fixture.detectChanges();

      compiled = fixture.nativeElement as HTMLElement;
      classes = compiled
        .querySelector('div')
        ?.attributes.getNamedItem('class')?.textContent;
      expect(classes).toContain('mat-app-background');
      expect(classes).toContain('dark-theme');
    });

    it('should ask for subscription if needed', () => {
      expect(push.requestSubscriptionOnDemand).toHaveBeenCalled();
      expect(push.updateNotificationsOnDemand).toHaveBeenCalled();
    });
  });
});
