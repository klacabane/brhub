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
      m('div[class="ui grid"]', [
        app.usermenu(),
        m('div[class="ten wide column"]', [
          m('button[type="button"][class="ui tiny button"][style="margin-right: 5px;"]', {onclick: module.vm.submit}, 'new item'),
          m('button[type="button"][class="ui tiny button"]', {onclick: module.vm.newTheme}, 'new theme'),
          app.grid.view(module.vm.grid)
        ])
      ])
    ]
  }

  return module;
}
