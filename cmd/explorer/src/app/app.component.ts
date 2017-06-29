import { Component } from '@angular/core';
import {Router, NavigationStart, NavigationEnd, NavigationCancel, NavigationError, Event as RouterEvent} from "@angular/router";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {

  private loading:boolean;
  constructor(private router: Router) {

    router.events.subscribe((event: RouterEvent) => {
      this.navigationInterceptor(event);
    });
  }


  // Shows and hides the loading spinner during RouterEvent changes
  navigationInterceptor(event: RouterEvent): void {
    if (event instanceof NavigationStart) {
      console.log("Navigation -start");
      this.loading = true;
    }
    if (event instanceof NavigationEnd) {
      console.log("Navigation -end");
      this.loading = false;
    }

    // Set loading state to false in both of the below events to hide the spinner in case a request fails
    if (event instanceof NavigationCancel) {
      console.log("Navigation -canceled");
      this.loading = false;
    }
    if (event instanceof NavigationError) {
      console.log("Navigation -error");
      this.loading = false;
    }
  }
}
