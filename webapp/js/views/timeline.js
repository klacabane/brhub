var app = app || {};

(function(module) {
  'use strict';

  module.view = function(ctrl) {
    return m("div", [
      m("button", {onclick: ctrl.signout}, "Signout"),
      m("button", {onclick: ctrl.newB}, "New B"),
      m("button", {onclick: ctrl.newItem}, "New Item"),
      app.grid.view(ctrl.grid)
    ]);
  }
 })(app.timeline = app.timeline || {});
