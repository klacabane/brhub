var app = app || {};

(function(module) {
  'use strict';

  module.controller = function() {
    if (storage.getUser() === null) return m.route('/');
  }

  module.view = function(ctrl) {
    return [
      app.banner(),
      app.usermenu(0)
    ];
  }
})(app.profile = {});
