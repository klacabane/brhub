var app = app || {};

(function(module) {
  'use strict';

  var vm = {}

  module.controller = function() {
    var user = storage.getUser();
    if (user === null) return m.route('/');

    this.grid = new app.grid.controller({src: 'timeline'});

    this.newItem = function() {
      User.newItem()
        .then(function(item) {console.log(item)}, app.utils.processError);
    };

    this.signout = function() {
      storage.setUser(null);
      m.route('/');
    };
  }
})(app.timeline = app.timeline || {});
