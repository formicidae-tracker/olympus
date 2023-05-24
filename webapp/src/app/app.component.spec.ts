import { TestBed, ComponentFixture } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { CoreModule } from './core/core.module';
import { NavToolbarComponent } from './nav-toolbar/nav-toolbar.component';
import { UserSettingsService } from './core/user-settings.service';

describe('AppComponent', () => {
  let service: UserSettingsService;

  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      imports: [RouterTestingModule, CoreModule],
      declarations: [AppComponent, NavToolbarComponent],
    });

    service = TestBed.inject(UserSettingsService);
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
  });
});
