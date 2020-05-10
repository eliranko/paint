import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit {
  @ViewChild('canvas', { static: true }) canvas: ElementRef<HTMLCanvasElement> | null = null;
  ctx: CanvasRenderingContext2D | null | undefined = null;
  isDrawing = false;
  x = 0;
  y = 0;

  constructor() { }

  ngOnInit(): void {
    this.ctx = this.canvas?.nativeElement.getContext('2d');
  }

  mouseDownCanvas(e: MouseEvent) {
    this.x = e.offsetX;
    this.y = e.offsetY;
    this.isDrawing = true;
  }

  mouseMoveCanvas(e: MouseEvent) {
    if (!this.isDrawing) return;
    this.drawLine(this.x, this.y, e.offsetX, e.offsetY);
    this.x = e.offsetX;
    this.y = e.offsetY;
  }

  mouseUpCanvas(e: MouseEvent) {
    if (!this.isDrawing) return;
    this.drawLine(this.x, this.y, e.offsetX, e.offsetY);
    this.x = 0;
    this.y = 0;
    this.isDrawing = false;
  }

  drawLine(x1: number, y1: number, x2: number, y2: number) {
    if (!this.ctx) return;

    this.ctx.beginPath();
    this.ctx.lineWidth = 3;
    this.ctx.strokeStyle = 'white';
    this.ctx.moveTo(x1, y1);
    this.ctx.lineTo(x2, y2);
    this.ctx.stroke();
    this.ctx.closePath();
  }
}
