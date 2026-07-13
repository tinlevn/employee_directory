import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatToolbarModule } from '@angular/material/toolbar';
import { HelpdialogComponent } from './helpdialog/helpdialog.component';
import { Observable } from 'rxjs';
import { ThemeService } from '../_services/theme.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatSlideToggleModule,
    MatToolbarModule
  ],
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


