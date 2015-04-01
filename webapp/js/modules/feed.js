var app = app || {};
app.feed = function(src) {
  var module = {};

  module.vm = {
    init: function() {
      this.user = m.prop(storage.getUser());
      if (this.user() === null) return m.route('/auth');

      this.banner = app.banner(this.user());
      this.grid = new app.grid.controller({src: src || m.route.param('name')});
    },
    signout: function() {
      storage.setUser(null);
      m.route('/auth');
    },
    submit: function() {
      m.route('/submit');
    },
    newBrhub: function() {
      m.route('/brhub');
    }
  }

  module.controller = function() {
    module.vm.init();
  }

  module.view = function() {
    return [
      module.vm.banner.view(),
      m('div[class="content"]', [
        m('button[type="button"][class="btn btn-default btn-sm"][style="margin-right: 5px;"]', {onclick: module.vm.submit}, 'new item'),
        m('button[type="button"][class="btn btn-default btn-sm"]', {onclick: module.vm.newBrhub}, 'new theme'),
        app.grid.view(module.vm.grid)
      ])
    ]
  }

  return module;
}

