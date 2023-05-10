import { TestBed, ComponentFixture } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { CoreModule } from './core/core.module';
import { NavToolbarComponent } from './nav-toolbar/nav-toolbar.component';
import { ThemeService } from './core/services/theme.service';

describe('AppComponent', () => {
  let service: ThemeService;

  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      imports: [RouterTestingModule, CoreModule],
      declarations: [AppComponent, NavToolbarComponent],
    });

    service = TestBed.inject(ThemeService);
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have a title that contains 'Olympus'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app.title).toContain('Olympus');
  });

  describe('render', () => {
    let compiled: HTMLElement;
    let fixture: ComponentFixture<AppComponent>;

    beforeEach(() => {
      fixture = TestBed.createComponent(AppComponent);
      fixture.detectChanges();
      compiled = fixture.nativeElement as HTMLElement;
    });
    it('should render title in navigation bar', () => {
      expect(
        compiled.querySelector('.mat-toolbar span')?.textContent
      ).toContain('FORT Olympus');
    });
    it('should apply mode accordingly', () => {
      let classes = compiled
        .querySelector('div')
        ?.attributes.getNamedItem('class')?.textContent;

      expect(classes).toContain('mat-app-background');
      expect(classes).not.toContain('dark-theme');

      service.setDarkTheme(true);
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
