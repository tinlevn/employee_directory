import { Component, OnInit } from '@angular/core';
import { MatDialogModule } from '@angular/material/dialog';

@Component({
  selector: 'app-helpdialog',
  standalone: true,
  imports: [MatDialogModule],
  templateUrl: './helpdialog.component.html'
})
export class HelpdialogComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

}
