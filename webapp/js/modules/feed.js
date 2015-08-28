var app = app || {};

app.feed = function(src) {
  var module = {};

  module.vm = {
    init: function() {
      this.user = m.prop(storage.getUser());
      if (this.user() === null) return m.route('/auth');

      this.grid = new app.grid.controller({src: src || m.route.param('name')});
    },
    signout: function() {
      storage.setUser(null);
      m.route('/auth');
    },
    submit: function() {
      m.route('/submit');
    },
    newTheme: function() {
      m.route('/theme');
    }
  }

  module.controller = function() {
    module.vm.init();
  }

  module.view = function() {
    return [
      app.banner(),
      m('div.ui.grid', [
        app.usermenu(1),
        m('div.ten.wide.column', [
          m('button.ui.tiny.button[type="button"][style="margin-right: 5px;"]', {onclick: module.vm.submit}, 'new item'),
          m('button.ui.tiny.button[type="button"]', {onclick: module.vm.newTheme}, 'new theme'),
          app.grid.view(module.vm.grid)
        ])
      ])
    ]
  }

  return module;
}
