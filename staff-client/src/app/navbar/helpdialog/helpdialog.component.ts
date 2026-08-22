import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-helpdialog',
  imports: [MatButtonModule, MatDialogModule],
  templateUrl: './helpdialog.component.html'
})
export class HelpdialogComponent {
  protected readonly dialogRef = inject(MatDialogRef<HelpdialogComponent>);
}
