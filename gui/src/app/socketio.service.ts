import { Injectable } from '@angular/core';
import { Socket } from 'ngx-socket-io';
import { Observable } from 'rxjs';
import { Canvas } from './models/Canvas';
import { switchMap, map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class SocketioService {
  canvas: Observable<Canvas>;

  constructor(private socket: Socket) {
    this.canvas = this.socket.fromEvent<string>("canvas").pipe(
      map(data => {
        let obj = JSON.parse(data);
        let canvas = new Canvas(obj.name, "");
        canvas.uuid = obj.uuid;
        return canvas;
      })
    );
  }

  listenToCanvasEvent(): Observable<Canvas> {
    return this.canvas;
  }
}
