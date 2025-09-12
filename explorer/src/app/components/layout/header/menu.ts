/**
 * Data for building the navigation menu. Used by the header and the footer.
 *
 * Each item could be a link to a section of the Skycoin website or a submenu.
 * The properties of the items are:
 * name: Name to show in the menu.
 * href: The URL the item opens (ignored if the item is a submenu).
 * active: If the item corresponds to the page that is currenly being shown (the explorer).
 * target: Value for the taget property of the <a> tag, if the item is not a submenu.
 * menu: if the item is a submenu, an array with more items.
 * open: if the item is a submenu, indicates if the submenu is open (must be set to false).
 */
export default [
  {
    name: 'Downloads',
    href: 'https://www.skycoin.com/downloads',
    active: false,
    target: '_blank',
    open: false,
  },
  {
    name: 'Ecosystem',
    open: false,
    active: false,
    menu: [
      {
        name: 'Overview',
        href: 'https://www.skycoin.com/ecosystem',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'Skywire',
        href: 'https://www.skycoin.com/skywire',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'Obelisk',
        href: 'https://www.skycoin.com/obelisk',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'Fiber',
        href: 'https://www.skycoin.com/fiber',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'CX',
        href: 'https://www.skycoin.com/cx',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'CXO',
        href: 'https://www.skycoin.com/cxo',
        active: false,
        target: '_blank',
        open: false,
      },
    ],
  },
  {
    name: 'Skyminer',
    href: 'https://www.skycoin.com/skyminer',
    active: false,
    target: '_blank',
    open: false,
  },
  {
    name: 'Blog',
    href: 'https://www.skycoin.com/blog/',
    active: false,
    target: '_blank',
    open: false,
  },
  {
    name: 'Other',
    open: false,
    active: true,
    menu: [
      {
        name: 'Gallery',
        href: 'https://www.skycoin.com/gallery',
        active: false,
        target: '_blank',
        open: false,
      },
      {
        name: 'Explorer',
        href: '/',
        active: true,
        target: '_self',
        open: false,
      },
      {
        name: 'Explorer API',
        href: 'https://explorer.skycoin.com/api.html',
        active: false,
        target: '_blank',
        open: false,
      },
    ],
  },
];
