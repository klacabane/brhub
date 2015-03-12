var app = app || {};

(function(module) {
  'use strict';

  module.controller = function() {
    if (storage.getUser() === null) return m.route('/');

    var id = m.route.param('id');
    if (!/^[0-9a-fA-F]{24}$/.test(id)) return m.route('/timeline');

    this.item = m.prop({});
    this.commentModules = [];
    this.reply = new reply({item: id});

    Item.get(id).
      then(function(item) {
        this.commentModules = item.comments.map(function(c) {
          return new comments(c, true);
        });
        this.item(item);
      }.bind(this), app.utils.processError);
  }

  module.view = function(ctrl) {
    return m('div', [
      m('h2', ctrl.item().title),
      m('p', ctrl.item().content),
      ctrl.reply.view(),
      m('div', ctrl.commentModules.map(function(mod) {
        return mod.view();
      }))
    ]);
  }
 })(app.items = {});
