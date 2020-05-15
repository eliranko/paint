import { Component, OnInit, OnDestroy } from '@angular/core';
import { CanvasService } from '../canvas.service';
import { Canvas } from '../models/Canvas';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SocketioService } from '../socketio.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.css']
})
export class SidenavComponent implements OnInit, OnDestroy {
  canvases: Canvas[] = [];
  currentCanvas: string = "";
  subscriptions: Subscription[] = [];

  constructor(private canvasService: CanvasService, private snackBar: MatSnackBar, private socketio: SocketioService) { }

  ngOnInit(): void {
    this.canvasService.getCanvases().subscribe(canvases => {
      this.canvases = canvases;
    }, () => {
      this.snackBar.open("Server is unavailable at the moment");
    });

    this.subscriptions.push(this.socketio.listenToCanvasEvent().subscribe(this.onCanvasUpdate.bind(this)));
  }

  ngOnDestroy() {
    for (let sub of this.subscriptions) sub.unsubscribe();
  }

  onCanvasUpdate(canvas: Canvas) {
    this.canvases.push(canvas);
  }

  onCanvasClick(uuid: string) {
    this.currentCanvas = uuid;
    this.canvasService.updateCanvas(uuid);
  }
}
