import { Component, OnInit } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-save-canvas',
  templateUrl: './save-canvas.component.html',
  styleUrls: ['./save-canvas.component.css']
})
export class SaveCanvasComponent implements OnInit {
  name: string = "";

  constructor() { }

  ngOnInit(): void {
  }
}
