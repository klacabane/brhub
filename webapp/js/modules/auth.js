var app = app || {};

app.banner = function(user) {
  return {
    vm: {
      signout: function() {
        storage.setUser(null);
        m.route('/auth');
      }
    },
    view: function() {
      var childs = [
        m('a', {href: user ? '/#/' : '/#/auth'}, [
          m('img', {src: 'images/salty.png', width: 90, height: 90})
        ])
      ];

      if (user) {
        childs.push(
          m('div', {style: {'float': 'right'}}, [
            m('a', {href: '/#/profile'}, user.name),
            m('button[type="button"][class="btn btn-link btn-xs"]', {onclick: this.vm.signout}, 'Sign out'),
          ])
        );
      }

      return m('div', {
        style: {
          width: '100%',
          background: '#18bc9c',
          height: '110px',
          paddingTop: '10px',
        },
      }, [
        m('div', {
          style: {
            width: '50%',
            margin: 'auto',
          },
        }, childs)
      ]);
    }
  };
};

(function(module) {
  module.controller = function() {
    if (storage.getUser() !== null) return m.route('/timeline');

    this.name = m.prop('');
    this.password = m.prop('');
    this.banner = app.banner();

    this.signin = function(e) {
      e.preventDefault();
      User.signin(this.name(), this.password())
        .then(function(user) {
          storage.setUser(user);
          m.route('/timeline');
        }, app.utils.processError);
    }.bind(this);
  }

  module.view = function(ctrl) {
    return [
      ctrl.banner.view(),
      m('div[class="content"]', [
        m('form[class="form-inline"]', [
          m('input[class="form-control"][name="username"][placeholder="Name"][type="text"]', {
            onchange: m.withAttr('value', ctrl.name)
          }),
          m('input[class="form-control"][name="password"][placeholder="Password"][type="password"]', {
            onchange: m.withAttr('value', ctrl.password)
          }),
          m('button[class="btn btn-default"][type="submit"]', {onclick: ctrl.signin}, 'Sign in')
        ])
      ])
    ]
  }
})(app.auth = {});
