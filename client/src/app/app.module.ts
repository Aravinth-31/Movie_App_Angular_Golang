import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {HttpClientModule} from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { FormsModule } from '@angular/forms';
import { MoviesComponent } from './movies/movies.component';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login/login.component';
import { AccountComponent } from './account/account.component';  

@NgModule({
  declarations: [
    AppComponent,
    MoviesComponent,
    LoginComponent,
    AccountComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    CommonModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
