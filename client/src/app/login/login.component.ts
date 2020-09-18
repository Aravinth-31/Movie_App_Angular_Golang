import { Component, OnInit } from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  email:string='';
  password:string='';
  error:boolean=false;
  lname:string='';
  lpass1:string='';
  lpass2:string='';
  lemail:string='';
  constructor(private http:HttpClient) { }
  signin(){
    this.http.post("http://localhost:3001/SignIn",JSON.stringify({email:this.email,password:this.password}))
    .subscribe(res=>{
      console.log(res);
      if(res==="success"){
        console.log(res)
        sessionStorage.setItem("userEmail",this.email);
        sessionStorage.setItem("logged",JSON.stringify(true));
        this.error=false;
        window.location.replace("http://localhost:4200/");
      }
      else
        this.error=true;
    })
  }
  signup(){
      this.http.post("http://localhost:3001/SignUp",JSON.stringify({name:this.lname,password:this.lpass1,email:this.lemail}))
      .subscribe(res=>{
        console.log(res);
        if(res=="success"){
          sessionStorage.setItem("userEmail",this.lemail);
          sessionStorage.setItem("logged",JSON.stringify(true));
          this.error=false;
          window.location.replace("http://localhost:4200/");
        }
        else  
          this.error=true;
      })
  }

  ngOnInit(): void {
  }

}
