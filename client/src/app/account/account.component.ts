import { Component, OnInit } from '@angular/core';
import '../movies/Reserve';
import { Reserve } from '../movies/Reserve';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.css']
})
export class AccountComponent implements OnInit {
  user:string='';
  reserves:any=[];
  flag:boolean=false;
  constructor(private http:HttpClient) { 
    this.user=sessionStorage.getItem("userEmail");
    this.reserves=JSON.parse(sessionStorage.getItem('reserves'))||[];
    console.log(this.user);
    http.post("http://localhost:3001/getTickets",JSON.stringify({email:this.user}))
    .subscribe(data=>{
      console.log(data)
      this.reserves=data;
      this.flag=true;
    });
  }
  ngOnInit();
  ngOnInit(): void {
  }

}
