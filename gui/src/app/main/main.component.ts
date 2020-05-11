import { Component, OnInit, ViewChild, ElementRef, OnDestroy } from '@angular/core';
import { fabric } from 'fabric';
import { CanvasService } from '../canvas.service';
import { MatDialog } from '@angular/material/dialog';
import { SaveCanvasComponent } from '../save-canvas/save-canvas.component';
import { Canvas } from '../models/Canvas';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit, OnDestroy {
  canvas: any;
  intervalId: any;
  storageKey = "canvas";
  subscriptions: Subscription[] = [];

  constructor(private canvasService: CanvasService,
    public dialog: MatDialog) { }

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

    this.subscriptions.push(this.canvasService.currentCanvas.subscribe(this.listenCanvasUpdates.bind(this)));
  }

  ngOnDestroy() {
    clearInterval(this.intervalId);
  }

  listenCanvasUpdates(uuid: string) {
    if (!uuid) return;

    this.canvasService.getCanvas(uuid).subscribe(canvas => this.canvas.loadFromJSON(canvas.data));
  }

  persistCanvas() {
    localStorage.setItem(this.storageKey, this.serializeCanvas());
  }

  save() {
    let dialog = this.dialog.open(SaveCanvasComponent, {
      width: '400px',
    });
    dialog.afterClosed().subscribe(name => { //unsubscribe?
      if (!name) return;

      this.canvasService.postCanvas(new Canvas(name, this.serializeCanvas())).subscribe();
      this.clear();
    });
  }

  clear() {
    this.canvas.clear();
  }

  serializeCanvas(): string {
    return JSON.stringify(this.canvas);
  }
}
