var app = app || {};

(function(module) {
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

  module.view = function(ctrl) {
    return m("div", [
      m("div", [
        m("input[name='username'][placeholder='Name'][type='text']", {
          onchange: m.withAttr("value", ctrl.name)
        }),
        m("input[name='password'][placeholder='Password'][type='password']", {
          onchange: m.withAttr("value", ctrl.password)
        }),
        m("input[type='submit'][value='swag me in']", {onclick: ctrl.signin})
      ])
    ]);
  }
})(app.auth = {});
