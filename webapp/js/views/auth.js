var app = app || {};

(function(module) {
  'use strict';

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
})(app.auth = app.auth || {});
