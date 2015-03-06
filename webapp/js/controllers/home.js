var app = app || {};

(function (module) {
  'use strict';

  var vm = {
    redirect: function (user) {
      storage.setUser(user);
      m.route('/dashboard');
    }
  }

  module.controller = function () {
    if (storage.getUser() !== null) return m.route('/dashboard');

    this.name = m.prop('');
    this.password = m.prop('');

    this.signin = function () {
      User.signin(this.name(), this.password())
        .then(vm.redirect, app.utils.processError);
    }.bind(this);
  }
})(app.home = app.home || {});
