import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { HelpdialogComponent } from './helpdialog/helpdialog.component';
import { Observable } from 'rxjs';
import { ThemeService } from '../_services/theme.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  isDarkMode$: Observable<boolean>;
  constructor(public dialog: MatDialog, private themeService: ThemeService) {
    this.isDarkMode$ = this.themeService.darkMode$;
  }
  ngOnInit(): void {
    this.themeService.initializeTheme();
  }

  openHelp() {
    this.dialog.open(HelpdialogComponent);
  }
  toggleDarkMode() {
    this.themeService.toggleDarkMode();
  }
}


