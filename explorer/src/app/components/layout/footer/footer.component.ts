import { Component } from '@angular/core';

import menu from '../header/menu';

/**
 * Skycoin footer. To activate it, FooterConfig.useGenericFooter (in app.config.ts)
 * must be false. Read the docs for more information.
 */
@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {

  /**
   * Data for building the contact icon list.
   */
  links = [
    /*
    {
        href: 'https://medium.com/@Skycoinproject',
        img: 'medium.svg',
    },
    {
        href: 'https://twitter.com/skycoinproject',
        img: 'twitter.svg',
    },
    {
        href: 'https://www.facebook.com/SkycoinOfficial',
        img: 'facebook.svg',
    },
    {
        href: 'https://www.instagram.com/skycoinproject/',
        img: 'instagram.svg',
    },
    */
    {
        href: 'https://github.com/skycoin/skycoin',
        img: 'github.svg',
    },
    /*
    {
        href: 'https://www.youtube.com/c/Skycoin',
        img: 'youtube.svg',
    },
    {
        href: 'https://www.reddit.com/r/skycoin',
        img: 'reddit.svg',
    },
    {
        href: 'https://itunes.apple.com/nl/podcast/skycoin/id1348472259?l=en',
        img: 'apple.svg',
    },
    {
        href: 'https://discord.gg/EgBenrW',
        img: 'discord.svg',
    },
    */
    {
        href: 'https://t.me/Skycoin',
        img: 'telegram.svg',
    },
    /*
    {
        href: 'https://www.linkedin.com/company/skycoin/',
        img: 'linkedin.svg',
    }
    */
  ];

  /**
   * Returns the data for building the navigation menu.
   */
  getMenu() {
    return menu;
  }

}
