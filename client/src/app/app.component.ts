import { Component } from '@angular/core';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http'
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  islogged:boolean=false;
  constructor(private router:Router,private http:HttpClient){
    this.islogged=JSON.parse(sessionStorage.getItem('logged'))||false;
    // http.get("http://localhost:3001/hello-world")
    // .subscribe(data=>{
    //   console.log(data)
    // })
  }

  redirect() {
    console.log("movies");
    this.router.navigate(['./movies']);
  }
  redirectLogin() {
    console.log("login");
    this.router.navigate(['./login']);
  }
  redirectAccount(){
    console.log("account");
    this.router.navigate(['./account']);
  }
  logOut(){
    sessionStorage.setItem("logged",JSON.stringify(false));
    sessionStorage.setItem("userEmail",JSON.stringify(""));
    window.location.replace("http://localhost:4200/movies");
  }
}
