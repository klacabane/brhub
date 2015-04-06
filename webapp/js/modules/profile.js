var app = app || {};

(function(module) {
  'use strict';

  module.controller = function() {
    var user = storage.getUser();
    if (user === null) return m.route('/');
  }

  module.view = function(ctrl) {
    return [
      app.banner()
    ];
  }
})(app.profile = {});
