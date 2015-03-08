var app = app || {};

(function(module) {
  'use strict';

  module.controller = function() {
    if (storage.getUser() !== null) return m.route('/timeline');

    this.name = m.prop('');
    this.password = m.prop('');

    this.signin = function() {
      User.signin(this.name(), this.password())
        .then(function(user) {
          storage.setUser(user);
          m.route('/timeline');
        }, app.utils.processError);
    }.bind(this);
  }
})(app.auth = app.auth || {});
