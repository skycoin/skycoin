import { Component } from '@angular/core';
import { Router, NavigationEnd, Event as RouterEvent } from '@angular/router';

import { HeaderConfig, FooterConfig } from 'app/app.config';
import { LanguageService } from 'app/services/language/language.service';
import { ExplorerService } from 'app/services/explorer/explorer.service';
import { ApiService } from './services/api/api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  // Get the settings for the generic header and the generic footer.
  headerConfig = HeaderConfig;
  footerConfig = FooterConfig;

  constructor(router: Router, languageService: LanguageService, explorerService: ExplorerService, apiService: ApiService) {
    // Initialize the services. The api service must be initialized first.
    apiService.initialize();
    languageService.initialize();
    explorerService.initialize();

    // Prevent automatic scroll retoration during navigation.
    if (history.scrollRestoration) {
      history.scrollRestoration = 'manual';
    }
    /*
    // Return the scroll to the top after navigating.
    router.events.subscribe((event: RouterEvent) => {
      if (event instanceof NavigationEnd) {
        window.scrollTo(0, 0);
      }
    });
    */
  }
}
