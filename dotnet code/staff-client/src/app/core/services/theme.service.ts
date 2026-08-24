import { Injectable, effect, signal } from '@angular/core';

const DARK_MODE_STORAGE_KEY = 'staff-search.darkMode';
const DARK_THEME_CLASS = 'dark-theme';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  readonly isDarkMode = signal(this.readStoredPreference());

  constructor() {
    effect(() => {
      const dark = this.isDarkMode();
      document.body.classList.toggle(DARK_THEME_CLASS, dark);
      localStorage.setItem(DARK_MODE_STORAGE_KEY, JSON.stringify(dark));
    });
  }

  toggleDarkMode(): void {
    this.isDarkMode.update((dark) => !dark);
  }

  private readStoredPreference(): boolean {
    try {
      return JSON.parse(localStorage.getItem(DARK_MODE_STORAGE_KEY) ?? 'false') === true;
    } catch {
      return false;
    }
  }
}
