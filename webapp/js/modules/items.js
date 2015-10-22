var app = app || {};

(function(module) {
  'use strict';

  var vm = {
    item: m.prop({}),
    commentModules: [],
    reply: m.prop(''),
    sendReply: function(e) {
      e.preventDefault();

      if (!vm.reply().length) return;

      Comments.create({
        item: vm.item().id,
        content: vm.reply()
      }).then(function(comment) {
        vm.commentModules.splice(0, 0, new commentModule(comment));
        vm.reply('');
      }, app.utils.processError);
    }
  }

  module.controller = function() {
    var user = storage.getUser();
    if (user === null) return m.route('/');

    var id = m.route.param('id');
    if (!/^[0-9a-fA-F]{24}$/.test(id)) return m.route('/timeline');

    Item.get(id).then(function(item) {
      vm.item(item);
      vm.commentModules = item.comments
        .map(function(comment) {
          return new commentModule(comment);
        });
    }.bind(this), app.utils.processError);
  }

  module.view = function(ctrl) {
    return [
      app.banner(),
      m('div[class="ui grid"]', [
        app.usermenu(),
        m('div[class="ten wide column"]', [
          m('a[class="btn-link"]', {href: vm.item().link}, vm.item().title),
          m('p', vm.item().content),
          m('form[class="ui reply form comment-form"]', [
              m('textarea', {
                onchange: m.withAttr('value', vm.reply),
                value: vm.reply()
              }),
            m('button[class="ui mini blue button"]', {onclick: vm.sendReply}, 'Reply')
          ]),
          m('div[class="ui threaded comments"]', 
            vm.commentModules.map(function(mod) {
              return mod.view();
            })
          )
        ])
      ])
    ]
  }
 })(app.items = {});
