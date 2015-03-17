var app = app || {};

(function(module) {
  'use strict';

  module.controller = function() {
    if (storage.getUser() === null) return m.route('/');

    var id = m.route.param('id');
    if (!/^[0-9a-fA-F]{24}$/.test(id)) return m.route('/timeline');

    this.item = m.prop({});
    this.commentModule;

    Item.get(id).
      then(function(item) {
        this.commentModule = new commentModule(item);
        this.item(item);
      }.bind(this), app.utils.processError);
  }

  module.view = function(ctrl) {
    return m('div', [
      m('h2', ctrl.item().title),
      m('p', ctrl.item().content),
      ctrl.commentModule.view(),
    ]);
  }
 })(app.items = {});
