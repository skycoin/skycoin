import {Component, OnInit, Output, EventEmitter} from '@angular/core';
import {SkycoinBlockchainPaginationService} from "./skycoin-blockchain-pagination.service";

@Component({
  selector: 'app-skycoin-pagination',
  templateUrl: './skycoin-pagination.component.html',
  styleUrls: ['./skycoin-pagination.component.css']
})
export class SkycoinPaginationComponent implements OnInit {

  @Output() onChangePage = new EventEmitter<any>();

  private numberOfBlocks:number;

  private currentPage:number;

  private pagesToShowAtATime:number;

  private pages:any;

  private pageStartPointer:number;

  private currentPages:number[];
  private pageEndPointer:number;

  private noUpcoming:boolean;

  constructor(private paginationService:SkycoinBlockchainPaginationService) {
    this.numberOfBlocks =0;
    this.currentPage = 1;
    this.currentPages=[];
    this.pagesToShowAtATime=5;
    this.pageStartPointer = this.currentPage;
    this.pageEndPointer = this.currentPage;
    this.noUpcoming = false;
  }

  ngOnInit() {
    this.paginationService.fetchNumberOfBlocks().subscribe((numberOfBlocks)=>{
      this.numberOfBlocks = numberOfBlocks;
      this.onChangePage.emit([1,  this.numberOfBlocks]);
      this.pagesToShowAtATime = this.pagesToShowAtATime<numberOfBlocks?this.pagesToShowAtATime:this.numberOfBlocks;

      this.currentPages = [];
      for (var i = this.currentPage; i < this.currentPage+this.pagesToShowAtATime; i++) {
        this.currentPages.push(i);
      }

    })
  }


  changePage(pageNumber:any){
    this.onChangePage.emit([pageNumber, this.numberOfBlocks]);
    this.currentPage = pageNumber;
    return false;
  }

  loadUpcoming():boolean{
    if(this.currentPages[0]*10+this.pagesToShowAtATime*10>=this.numberOfBlocks){
      this.noUpcoming = true;
      return false;
    }
    this.onChangePage.emit([this.currentPages[0]+this.pagesToShowAtATime,  this.numberOfBlocks]);
    this.currentPage = this.currentPages[0]+this.pagesToShowAtATime;

    this.currentPages = [];
    for (var i = this.currentPage; i < this.currentPage+this.pagesToShowAtATime && i<=this.numberOfBlocks; i++) {
      if(i*10-this.numberOfBlocks< 10){
        this.currentPages.push(i);
      }
      else{
        this.noUpcoming = true;
      }

    }


    return false;
  }

  loadPrevious():boolean{
    this.noUpcoming = false;
    if(this.currentPages[0]<=1){
      return false;
    }

    if(this.currentPages[0]-this.pagesToShowAtATime<=0){
      this.currentPages = [];
      this.currentPage = 1;
      this.onChangePage.emit([1, this.numberOfBlocks]);
      for (var i = this.currentPage; i < this.currentPage+this.pagesToShowAtATime; i++) {
        this.currentPages.push(i);
      }

    }
    else{

      this.onChangePage.emit([this.currentPages[0]-this.pagesToShowAtATime, this.numberOfBlocks]);
      this.currentPage = this.currentPages[0]-this.pagesToShowAtATime;

      this.currentPages = [];
      for (var i = this.currentPage; i <this.currentPage+this.pagesToShowAtATime; i++) {
        if(i*10<=this.numberOfBlocks){
          this.currentPages.push(i);
        }

      }
    }


    return false;
  }

}
