var app = app || {};

app.banner = function() {
  return m('div.banner');
};

app.usermenu = function() {
  var user = storage.getUser();
  return m('div.four.wide.column', [
    m('div.ui.secondary.vertical.menu', [
      m('a.item', {href: '/#/users/'+user.name}, user.name),
      m('a.item', {href: '/#/'}, 'home'),
      m('a.item', {onclick: function() {
        storage.setUser(null);
        m.route('/auth');
      }},'signout'),
      m('div.item', [
        m('div.ui.fluid.category.search', [
          m('div.ui.mini.icon.input', [
            m('input.prompt[type="text"][placeholder="search"]', {
              config: function(elem, init, ctx) {
                if (!init) {
                  $('.ui.search').search({
                    type: 'category',
                    apiSettings: {
                      url: '/api/search/{$query}',
                      beforeXHR: function(xhr) {
                        xhr.setRequestHeader('X-token', user.token);
                      }
                    }
                  });
                }
              }
            }),
            m('i.search.icon')
          ]),
          m('div.results')
        ])
      ])
    ])
  ]);
}

var auth = app.auth = {};
auth.controller = function() {
  if (storage.getUser() !== null) return m.route('/timeline');

  this.name = m.prop('');
  this.password = m.prop('');
  this.signin = function(e) {
    e.preventDefault();

    User.signin(this.name(), this.password())
      .then(function(user) {
        storage.setUser(user);
        m.route('/timeline');
      }, app.utils.processError);
  }.bind(this);
}

auth.view = function(ctrl) {
  return [
    app.banner(),
    m('div.wrapper', [
      m('form.ui.form', [
        m('div.three.fields', [
          m('div.field', [
            m('input[name="username"][placeholder="Name"][type="text"]', {
              onchange: m.withAttr('value', ctrl.name)
            }),
          ]),
          m('div.field', [
            m('input[name="password"][placeholder="Password"][type="password"]', {
              onchange: m.withAttr('value', ctrl.password)
            }),
          ]),
          m('div.field', [
            m('button[class="ui button"][type="submit"]', {onclick: ctrl.signin}, 'Sign in')
          ])
        ])
      ])
    ])
  ]
}
