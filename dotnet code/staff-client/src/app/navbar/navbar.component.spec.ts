import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';

import { ThemeService } from '../core/services/theme.service';
import { HelpdialogComponent } from './helpdialog/helpdialog.component';
import { NavbarComponent } from './navbar.component';

describe('NavbarComponent', () => {
  let fixture: ComponentFixture<NavbarComponent>;
  let component: NavbarComponent;
  let themeService: ThemeService;

  beforeEach(async () => {
    localStorage.clear();
    document.body.classList.remove('dark-theme');

    await TestBed.configureTestingModule({
      imports: [NavbarComponent],
      providers: [provideZonelessChangeDetection()]
    }).compileComponents();

    themeService = TestBed.inject(ThemeService);
    fixture = TestBed.createComponent(NavbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should toggle the theme through the ThemeService', () => {
    expect(themeService.isDarkMode()).toBeFalse();

    const toggle: HTMLElement = fixture.nativeElement.querySelector('.theme-toggle-wrapper');
    toggle.click();
    fixture.detectChanges();

    expect(themeService.isDarkMode()).toBeTrue();
  });

  it('should reflect the current theme in the label', () => {
    const label: HTMLElement = fixture.nativeElement.querySelector('.theme-label');
    expect(label.textContent).toContain('Light');

    themeService.toggleDarkMode();
    fixture.detectChanges();

    expect(label.textContent).toContain('Dark');
  });

  it('should open the help dialog', () => {
    // Material 22 declares MatDialog with @Service, so instance identity is not
    // shared with the TestBed injector — spy on the prototype instead.
    const openSpy = spyOn(MatDialog.prototype, 'open');

    component.openHelp();

    expect(openSpy).toHaveBeenCalledWith(HelpdialogComponent);
  });
});
