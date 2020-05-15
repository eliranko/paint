import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MainViewEncapsulatorComponent } from './main-view-encapsulator/main-view-encapsulator.component';
import { NotFoundComponent } from './not-found/not-found.component';


const routes: Routes = [
  { path: 'draw', component: MainViewEncapsulatorComponent },
  { path: '', redirectTo: '/draw', pathMatch: 'full' },
  { path: '**', component: NotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
