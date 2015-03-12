var app = app || {};
app.feed = function(src) {
  var module = {};

  module.vm = {
    init: function() {
      if (storage.getUser() === null) return m.route('/auth');

      this.grid = new app.grid.controller({src: src || m.route.param('name')});
    },
    signout: function() {
      storage.setUser(null);
      m.route('/auth');
    },
    submit: function() {
      m.route('/submit');
    }
  };

  module.controller = function() {
    module.vm.init();
  }

  module.view = function(ctrl) {
    return m('div', [
      m('button', {onclick: module.vm.signout}, 'Signout'),
      m('button', {onclick: module.vm.submit}, 'New'),
      app.grid.view(module.vm.grid)
    ]);
  };

  return module;
}

