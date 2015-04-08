var app = app || {};

app.banner = function() {
  return m('div[class="banner"]');
};

app.usermenu = function() {
  var user = storage.getUser();
  return m('div[class="four wide column"]', [
    m('div[class="ui secondary vertical menu"]', [
      m('a[class="item"]', {href: '/#/users/'+user.name}, user.name),
      m('a[class="item"]', {href: '/#/'}, 'home'),
      m('a[class="item"]', {onclick: function() {
        storage.setUser(null);
        m.route('/auth');
      }},'signout'),
      m('div[class="item"]', [
        m('div[class="ui fluid category search"]', [
          m('div[class="ui mini icon input"]', [
            m('input[class="prompt"][type="text"][placeholder="search"]', {onfocus: function() {
              $('.ui.search').search({
                type: 'category',
                apiSettings: {
                  url: '/api/search/{$query}',
                  beforeXHR: function(xhr) {
                    xhr.setRequestHeader('X-token', user.token);
                  }
                }
              });
            }}),
            m('i[class="search icon"]')
          ]),
          m('div[class="results"]')
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
    m('div[class="wrapper"]', [
      m('form[class="ui form"]', [
        m('div[class="three fields"]', [
          m('div[class="field"]', [
            m('input[name="username"][placeholder="Name"][type="text"]', {
              onchange: m.withAttr('value', ctrl.name)
            }),
          ]),
          m('div[class="field"]', [
            m('input[name="password"][placeholder="Password"][type="password"]', {
              onchange: m.withAttr('value', ctrl.password)
            }),
          ]),
          m('div[class="field"]', [
            m('button[class="ui button"][type="submit"]', {onclick: ctrl.signin}, 'Sign in')
          ])
        ])
      ])
    ])
  ]
}
