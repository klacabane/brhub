var app = app || {};

(function (module) {
  'use strict';

  module.view = function (ctrl) {
    console.log(ctrl.item())
  }

})(app.items = app.items || {})
