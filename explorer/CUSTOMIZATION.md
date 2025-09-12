# Customizing the explorer

The explorer is designed in such a way that it is possible to make simple customizations with relative ease.

## Page title

To change the page title (the one shown in the browser tab), simply change the content of the `<title>` tag in [index.html](src/index.html).

## Favicon

The favicon file is [favicon.ico](src/favicon.ico), so you can change the favicon by replacing that file. You can also change the content of the `<link>` tag used for setting the favicon in [index.html](src/index.html).

## Basic identifiers

The explorer automatically tries to obtain from the node the name of the coins and hours. If, for example, the explorer is connected to a Skycoin node, it will automatically get the "SKY" name from the node and use it to identify the amounts of coins.

However, because old versions of the node do not return that information, you should add the coin identifiers in the `CoinIdentifiers` object inside the [app.config.ts](src/app/app.config.ts) file.

## QR codes

To simplify copying addresses with mobile devices, the explorer shows qr codes. The qr codes contain the addresses prefixed with the name of the coin, in this way: `skycoin:abcd...` (More information in this link: https://github.com/bitcoin/bips/blob/master/bip-0021.mediawiki). If needed, the `skycoin:` prefix can be changed by modifying the value of `QrConfig.prefix`, inside [app.config.ts](src/app/app.config.ts).

## Search functionality

The explorer is integrated with the browser's search functionality. This means that the user can search for elements of the blockchain directly from the search functionality of the browser. For this to work properly, you must modify the [search.xml](src/search.xml) file. Simply replace the `https://explorer.skycoin.com/` prefix of all URLs with the URL where the browser will reside (it is not necessary to replace the text that may be to the right of that prefix). Also, replace the `ShortName` and `Description` as deemed necessary.

You can find more information about the file format in [opensearch-1-1-draft-6.md](https://github.com/dewitt/opensearch/blob/master/opensearch-1-1-draft-6.md)

## Colors and general appearance

Most explorer colors and other general UI parameters are defined in [_variables.scss](src/assets/scss/_variables.scss). Making changes to that file is the quickest way to change the explorer appearance.

## Header

It is possible to completely customize the explorer header by manually modifying the code of [HeaderComponent](src/app/components/layout/header/header.component.ts). However, the explorer also has an easily configurable generic header.

The generic header may be activated simply by opening [app.config.ts](src/app/app.config.ts) and setting `HeaderConfig.useGenericHeader` to `true`. The generic header looks like this:

![GitHub Logo](/doc_images/header1.png)

### Customizing the generic header

The "website" link opens the URL specified in `HeaderConfig.genericHeaderUrl`, inside [app.config.ts](src/app/app.config.ts).

The logo may be modified by replacing the [logo.png]((src/assets/img/logo.png) file. In case of modifying the logo, it is important to keep in mind that the new image should be less than 110px high and 220px wide.

The header colors may be modified by changing the values of `$col-generic-header-background`, `$col-generic-header-text` and `$col-generic-header-border`, in [_variables.scss](src/assets/scss/_variables.scss). For example, consider the following values:
```
$col-generic-header-background: lightyellow;
$col-generic-header-text: darkseagreen;
$col-generic-header-border: lightpink;
```
The previous values generate the following style:

![GitHub Logo](/doc_images/header2.png)

## Footer

Similarly to what happens with the header, it is possible to completely customize the appearance of the explorer footer by manually modifying the code of [FooterComponent](src/app/components/layout/footer/footer.component.ts), as well as showing an easily configurable generic footer.

The generic footer may be activated simply by opening [app.config.ts](src/app/app.config.ts) and setting `FooterConfig.useGenericFooter` to `true`. The generic footer looks like this:

![GitHub Logo](/doc_images/footer1.png)

### Customizing the generic footer

The contact icons shown in the footer are defined in the `FooterConfig.contactLinks` array, inside [app.config.ts](src/app/app.config.ts). It is simply an array of objects, like the following:
```
contactLinks: [
  {
    url: "https://github.com/skycoin",
    content: "<i class='fab fa-github'></i>",
  },{
    url: "https://t.me/Skycoin",
    content: "<i class='fab fa-telegram'></i>",
  }
]
```
Each element has two properties:

- url: URL that will be opened when the user clicks on the icon.
- content: HTML code of the icon that will be displayed. The icon must fit within a 30px by 30px square.

Although it may seem that the second parameter is complicated to configure, its use is very simple, as the explorer uses Font Awesome v5. To add an icon simply open https://fontawesome.com/icons?d=gallery&m=free, click on the desired icon and, on the next page, copy the HTML code. For example, selecting the Twitter icon opens https://fontawesome.com/icons/twitter?style=brands, where it is possible to find the HTML code that must be placed in the `content` property, as shown in the following image (marked in red):

![GitHub Logo](/doc_images/code.png)

Taking all this into account, it is possible to add a link to a Twitter account simply by adding an element like the following one to the `FooterConfig.contactLinks` array:
```
{
  url: "https://twitter.com/skycoinproject",
  content: "<i class='fab fa-twitter'></i>",
}
```
Also, as the `content` property receives HTML code, it is possible to include any valid HTML element, its use is not limited to Font Awesome.

Additionally, the footer colors may be modified by changing the values of `$col-generic-footer-background`, `$col-generic-footer-text`, `$col-generic-footer-top-border` and `$col-generic-footer-middle-border` in [_variables.scss](src/assets/scss/_variables.scss). For example, consider the following values:
```
$col-generic-footer-background: lightgray;
$col-generic-footer-text: dodgerblue;
$col-generic-footer-top-border: #000000;
$col-generic-footer-middle-border: orange;
```
The previous values generate the following style:

![GitHub Logo](/doc_images/footer2.png)
