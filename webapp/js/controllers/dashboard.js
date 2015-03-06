var app = app || {};

(function (module) {
  'use strict';

  var vm = {}

  module.controller = function () {
    var user = storage.getUser();
    var grid = new app.dashboard.controller();

    if (user === null) return m.route('/');

    this.signout = function () {
      storage.setUser(null);
      m.route('/');
    };
  }
})(app.dashboard = app.dashboard || {});
