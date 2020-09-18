import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MoviesComponent } from './movies/movies.component';
import { LoginComponent } from './login/login.component';  
import { AccountComponent } from './account/account.component';  

const routes: Routes = [
  {path:'',redirectTo:'/movies',pathMatch:'full'},
  {path:'movies',component:MoviesComponent},
  {path:'login',component:LoginComponent},
  {path:'account',component:AccountComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
