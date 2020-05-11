import { Component, OnInit } from '@angular/core';
import { CanvasService } from '../canvas.service';
import { Canvas } from '../models/Canvas';

@Component({
  selector: 'app-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.css']
})
export class SidenavComponent implements OnInit {
  canvases: Canvas[] = [];
  constructor(private canvasService: CanvasService) { }

  ngOnInit(): void {
    this.canvasService.getCanvases().subscribe(canvases => {
      this.canvases = canvases;
    });
  }

  onCanvasClick(uuid: string) {
    this.canvasService.updateCanvas(uuid);
  }
}
