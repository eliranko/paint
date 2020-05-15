import { Component, OnInit } from '@angular/core';
import { CanvasService } from '../canvas.service';
import { Canvas } from '../models/Canvas';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.css']
})
export class SidenavComponent implements OnInit {
  canvases: Canvas[] = [];
  currentCanvas: string = "";

  constructor(private canvasService: CanvasService, private snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.canvasService.getCanvases().subscribe(canvases => {
      this.canvases = canvases;
    },
      err => {
        this.snackBar.open("Server is unavailable at the moment");
      });
  }

  onCanvasClick(uuid: string) {
    this.currentCanvas = uuid;
    this.canvasService.updateCanvas(uuid);
  }
}
