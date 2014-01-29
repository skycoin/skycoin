/*
 * angular-qrcode v2.0.0
 * (c) 2013 Monospaced http://monospaced.com
 * License: MIT
 */

angular.module('monospaced.qrcode', [])
  .directive('qrcode', ['$timeout', '$window', function($timeout, $window){

    var canvas2D = !!$window.CanvasRenderingContext2D,
        levels = {
          'L': 'Low',
          'M': 'Medium',
          'Q': 'Quartile',
          'H': 'High'
        },
        draw = function(context, qr, modules, tile){
          for (var row = 0; row < modules; row++) {
            for (var col = 0; col < modules; col++) {
              var w = (Math.ceil((col + 1) * tile) - Math.floor(col * tile)),
                  h = (Math.ceil((row + 1) * tile) - Math.floor(row * tile));
              context.fillStyle = qr.isDark(row, col) ? '#000' : '#fff';
              context.fillRect(Math.round(col * tile), Math.round(row * tile), w, h);
            }
          }
        };

    return {
      restrict: 'E',
      template: '<canvas></canvas>',
      link: function(scope, element, attrs){

        var domElement = element[0],
            canvas = element.find('canvas')[0],
            version = Math.max(1, Math.min(parseInt(attrs.version, 10), 10)) || 4,
            correction = attrs.errorCorrectionLevel in levels ? attrs.errorCorrectionLevel : 'M',
            trim = /^\s+|\s+$/g,
            qr = qrcode(version, correction);

        qr.make();

        var modules = qr.getModuleCount(),
            size = parseInt(attrs.size, 10) || modules * 2,
            tile = size / modules,
            render = function(qr, text){
              qr.addData(text);
              qr.make();
              if (canvas2D) {
                draw(context, qr, modules, tile);
              } else {
                domElement.innerHTML = qr.createImgTag(tile, 0);
              }
            };

        if (canvas2D) {
          var context = canvas.getContext('2d');
          canvas.width = canvas.height = size;
        }

        attrs.$observe('data', function(value){
          if (!value) {
            return;
          }
          var text = value.replace(trim, ''),
              qr = qrcode(version, correction);
          render(qr, text);
        });
      }
    };
  }]);
