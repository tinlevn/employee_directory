import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';

import { HelpdialogComponent } from './helpdialog/helpdialog.component';
import { ThemeService } from '../core/services/theme.service';

@Component({
  selector: 'app-navbar',
  imports: [MatButtonModule, MatDialogModule, MatIconModule, MatToolbarModule],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent {
  private readonly themeService = inject(ThemeService);
  private readonly dialog = inject(MatDialog);

  protected readonly isDarkMode = this.themeService.isDarkMode;

  protected toggleDarkMode(): void {
    this.themeService.toggleDarkMode();
  }

  protected openHelp(): void {
    this.dialog.open(HelpdialogComponent);
  }
}
