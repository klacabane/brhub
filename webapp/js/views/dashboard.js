var app = app || {};

(function (module) {
  'use strict';

  module.view = function (ctrl) {
    return m("div", [
      m("button", {onclick: ctrl.signout}, "Signout"),
      app.grid.view(ctrl.grid)
    ]);
  }
 })(app.dashboard = app.dashboard || {});
