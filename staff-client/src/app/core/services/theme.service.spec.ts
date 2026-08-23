import { TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';

import { ThemeService } from './theme.service';

const STORAGE_KEY = 'staff-search.darkMode';

describe('ThemeService', () => {
  beforeEach(() => {
    localStorage.clear();
    document.body.classList.remove('dark-theme');

    TestBed.configureTestingModule({
      providers: [provideZonelessChangeDetection()]
    });
  });

  it('should default to light mode when nothing is stored', () => {
    const service = TestBed.inject(ThemeService);

    expect(service.isDarkMode()).toBeFalse();
    expect(document.body.classList.contains('dark-theme')).toBeFalse();
  });

  it('should toggle dark mode, update the DOM and persist the choice', () => {
    const service = TestBed.inject(ThemeService);

    service.toggleDarkMode();
    TestBed.tick();

    expect(service.isDarkMode()).toBeTrue();
    expect(document.body.classList.contains('dark-theme')).toBeTrue();
    expect(localStorage.getItem(STORAGE_KEY)).toBe('true');

    service.toggleDarkMode();
    TestBed.tick();

    expect(service.isDarkMode()).toBeFalse();
    expect(document.body.classList.contains('dark-theme')).toBeFalse();
    expect(localStorage.getItem(STORAGE_KEY)).toBe('false');
  });

  it('should restore a stored dark-mode preference on creation', () => {
    localStorage.setItem(STORAGE_KEY, 'true');

    const service = TestBed.inject(ThemeService);

    expect(service.isDarkMode()).toBeTrue();
  });

  it('should fall back to light mode when storage holds invalid JSON', () => {
    localStorage.setItem(STORAGE_KEY, '{not json');

    const service = TestBed.inject(ThemeService);

    expect(service.isDarkMode()).toBeFalse();
  });
});
