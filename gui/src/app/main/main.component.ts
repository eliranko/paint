import { Component, OnInit, ViewChild, ElementRef, OnDestroy } from '@angular/core';
import { fabric } from 'fabric';
import { CanvasService } from '../canvas.service';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit, OnDestroy {
  canvas: any;
  intervalId: any;
  storageKey = "canvas";

  constructor(private canvasService: CanvasService) { }

  ngOnInit(): void {
    this.canvas = new fabric.Canvas("canvas", {
      isDrawingMode: true,
    });
    fabric.Object.prototype.transparentCorners = false;
    this.canvas.freeDrawingBrush.color = "white";
    this.canvas.freeDrawingBrush.width = 1;
    this.canvas.renderAll();

    setTimeout(() => {
      this.canvas.loadFromJSON(localStorage.getItem(this.storageKey));
      this.intervalId = setInterval(() => this.persistCanvas(), 1000);
    }, 1);
  }

  ngOnDestroy() {
    clearInterval(this.intervalId);
  }

  persistCanvas() {
    localStorage.setItem(this.storageKey, this.serializeCanvas());
  }

  save() {
    this.canvasService.postCanvas(this.serializeCanvas()).subscribe();
  }

  clear() {
    this.canvas.clear();
  }

  serializeCanvas(): string {
    return JSON.stringify(this.canvas);
  }
}
