import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Canvas } from './models/Canvas';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CanvasService {
  private canvasUpdateSource = new BehaviorSubject('');
  public currentCanvas = this.canvasUpdateSource.asObservable();

  constructor(private http: HttpClient) { }

  updateCanvas(uuid: string) {
    this.canvasUpdateSource.next(uuid);
  }

  getCanvases(): Observable<Canvas[]> {
    return this.http.get<Canvas[]>("api/canvas").pipe(
      retry(3),
      catchError(this.handleError)
    );
  }

  getCanvas(uuid: string): Observable<Canvas> {
    return this.http.get<Canvas>("api/canvas/" + uuid).pipe(
      retry(3),
      catchError(this.handleError)
    );
  }

  postCanvas(canvas: Canvas): Observable<void> {
    return this.http.post<void>("api/canvas", canvas).pipe(
      retry(3),
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error.message);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong,
      console.error(
        `Backend returned code ${error.status}, ` +
        `body was: ${error.error}`);
    }
    // return an observable with a user-facing error message
    return throwError(
      'Something bad happened; please try again later.');
  };
}
