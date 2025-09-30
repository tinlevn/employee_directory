import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {

  private darkMode = new BehaviorSubject<boolean>(this.getStoredTheme());
  darkMode$ = this.darkMode.asObservable();

  private getStoredTheme(): boolean {
    const storedTheme = localStorage.getItem('darkMode');
    return storedTheme ? JSON.parse(storedTheme) : false;
  }

  toggleDarkMode() {
    const newValue = !this.darkMode.value;
    localStorage.setItem('darkMode', JSON.stringify(newValue));
    this.darkMode.next(newValue);
    
    if (newValue) {
      document.body.classList.add('dark-theme');
    } else {
      document.body.classList.remove('dark-theme');
    }
  }

  initializeTheme() {
    if (this.darkMode.value) {
      document.body.classList.add('dark-theme');
    }
  }
}
