# Multi-coin compatibility

This wallet is not limited to one coin, but is compatible with most coins based on the Skycoin technology. For a
coin to be compatible, it must have a public node updated to at least version 0.24 and, ideally, its REST api
should have not been modified (if the api is modified, the wallet could fail depending on whether any of the
required endpoints were modified).

## Adding a coin

To add a coin, you must provide the required data about the coin. Most data must be added in a new class, in a new
file in the [src/app/coins](src/app/coins) folder. The class must inherit from
[BaseCoin](src/app/coins/basecoin.ts), the base class that defines the data a coin must have (it is mandatory to
add data to the abstract properties, the rest are optional). The data defined
in [BaseCoin](src/app/coins/basecoin.ts) are:

- `id`: a unique identifier for the coin (different from zero).
- `nodeUrl`: the URL of the public node that will be used as backend for the operations involving the coin.
- `coinName`: the full name of the coin.
- `coinSymbol`: the short name of the coin.
- `hoursName`: the name of the hours produced by the coin.
- `priceTickerId`: the ID of the coin in Coin Paprika.
- `coinExplorer`: the URL of the blockchain explorer of the coin.
- `imageName`: name of the file with the image that will be displayed in the header. The image must be stored in
[src/assets/img/coins](src/assets/img/coins) and it's size must be 1280x720px.
- `gradientName`: name of the file with the gradient that will be displayed in front of the header background.
The image must be stored in [src/assets/img/coins](src/assets/img/coins) and should have the smallest possible
size. If you do not want to add a gradient, you can simply ignore this property (more information below).
- `iconName`: name of the file with the small icon that will be displayed inside the button for changing the coin,
in the header. The icon will be shown when the coin is selected, so it should look good above the `imageName` and
`gradientName` images. The image must be stored in [src/assets/img/coins](src/assets/img/coins) and it's size
must be 20x20px. To make it look good in the design, the image must be in a .png file with transparent background.
- `bigIconName`: name of the file with the icon that will be displayed in the coin selection screen. The image must
be stored in [src/assets/img/coins](src/assets/img/coins) and it's size must be 64x64px. To make it look good in
the design, the image must be in a .png file with transparent background.

For an example of how to create the class and what form the data should have, please see
[SkycoinCoin](src/app/coins/skycoin.coin.ts).

Once the class is created and the image files are added, it is necessary to add a new instance of the newly created
class in the [CoinService.coins](src/app/services/coin.service.ts) array. That must be done in the
`CoinService.loadAvailableCoins` function.

After that, you must add the host name of the coin node to the `defaultCoinHosts` array, inside
[electron/src/electron-main.js](electron/src/electron-main.js). For example, if the node is in
`https://node.skycoin.net`, you must add a new entry with `node.skycoin.net`.

To finish, add the node URL (`https://node.skycoin.net`, for example) in the `connect-src` section of the
`Content-Security-Policy` header defined in [docker/images/nginx.conf](docker/images/nginx.conf).

## Additional information about the header images

The image indicated in `imageName` is the one that will be shown in the upper part of the wallet, as the background
image for the balance. The header places the image in such a way that it is rescaled in accordance with the size of
the window, always maintaining its aspect ratio, so it is not possible to know in advance what will be the area that
the user will see.

Taking into account that there will be text on top of the background image, it will probably be necessary to add
a transparent color or gradient on top of it. Adding a plain color is simple, since it can be added directly to
the background bitmap using an image editor, but gradients and other effects could give problems.

Due to the way in which the background image is rescaled, a gradient added directly into the background bitmap could
look bad in some windows sizes, that is the reason why an additional bitmap can be added with the `gradientName`
property. If an image is added to that property, the image will be shown in front of the background (but behind the
text) and will be rescaled so that it always occupies the header area, which means that it will always be fully
visible and that it will not maintain its aspect ratio.

An example of this combination is that when you select Skycoin as the active coin, an image with various gray coins is
shown as the background image with blue gradient above it. The image of the coins is
[skycoin-header.jpg](src/assets/img/coins/skycoin-header.jpg) and the image with the gradient is 
[skycoin-gradient.png](src/assets/img/coins/skycoin-gradient.png).

If no value is added to `imageName` and `gradientName`, the header shows a generic background.

## Hiding the multi-coin options

If the [CoinService.coins](src/app/services/coin.service.ts) array has less than two elements, the coin selection
options are hidden from the UI. This means that it is posible to remove the coin selection options by adding only
one coin in the `CoinService.loadAvailableCoins` function.