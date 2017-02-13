import {Component, OnInit, Output, EventEmitter} from '@angular/core';
import {SkycoinBlockchainPaginationService} from "./skycoin-blockchain-pagination.service";

@Component({
  selector: 'app-skycoin-pagination',
  templateUrl: './skycoin-pagination.component.html',
  styleUrls: ['./skycoin-pagination.component.css']
})
export class SkycoinPaginationComponent implements OnInit {

  @Output() onChangePage = new EventEmitter<number>();

  private numberOfBlocks:number;

  private currentPage:number;

  private pagesToShowAtATime:number;

  private pages:any;

  private pageStartPointer:number;

  private currentPages:number[];
  private pageEndPointer:number;

  constructor(private paginationService:SkycoinBlockchainPaginationService) {
    this.numberOfBlocks =0;
    this.currentPage = 1;
    this.currentPages=[];
    this.pagesToShowAtATime=5;
    this.pageStartPointer = this.currentPage;
    this.pageEndPointer = this.currentPage;
  }

  ngOnInit() {
    this.paginationService.fetchNumberOfBlocks().subscribe((numberOfBlocks)=>{
      this.numberOfBlocks = numberOfBlocks;
      this.pagesToShowAtATime = this.pagesToShowAtATime<numberOfBlocks?this.pagesToShowAtATime:this.numberOfBlocks;

      this.currentPages = [];
      for (var i = this.currentPage; i <= this.currentPage+4; i++) {
        this.currentPages.push(i);
      }

    })
  }

  setPage(currentPage:any){
    if(!(currentPage in this.pages)){

    }
  }

  changePage(pageNumber:any){
    this.onChangePage.emit(pageNumber);
    this.currentPage = pageNumber;
    return false;
  }

  loadUpcoming():boolean{

    this.onChangePage.emit(this.currentPages[0]+this.pagesToShowAtATime);
    this.currentPage = this.currentPages[0]+this.pagesToShowAtATime;

    this.currentPages = [];
    for (var i = this.currentPage; i <= this.currentPage+4; i++) {
      if(this.numberOfBlocks-i*10>=0){
        this.currentPages.push(i);
      }
      else if(this.numberOfBlocks-i*10>=-10){
        this.currentPages.push(i);
      }
    }


    return false;
  }

  loadPrevious():boolean{
    if(this.currentPages[0]<=1){
      return false;
    }
    this.onChangePage.emit(this.currentPages[0]-this.pagesToShowAtATime);
    this.currentPage = this.currentPages[0]-this.pagesToShowAtATime;

    this.currentPages = [];
    for (var i = this.currentPage; i <= this.currentPage+4; i++) {
      if(i*10<=this.numberOfBlocks){
        this.currentPages.push(i);
      }

    }

    return false;
  }

}
