var app = app || {};

var commentModule = function(parent) {
  var module = {};

  module.vm = {
    reply: m.prop(''),
    showForm: m.prop(false),
    sendReply: function(e) {
      e.preventDefault();

      if (!module.vm.reply().length) return;

      Comments.create({
        item: parent.item,
        parent: parent.id,
        content: module.vm.reply()
      }).then(function(comment) {
        module.vm.childModules.splice(0, 0, new commentModule(comment));
        module.vm.reply('');
        module.vm.showForm(false);
      }, app.utils.processError);
    },
    childModules: parent.comments.map(function(child) {
      return new commentModule(child);
    })
  };

  module.view = function() {
    return m('div.comment', [
      m('a.avatar', [
        m('img', {src: 'images/avatar.jpg'})
      ]),
      m('div.content', [
        m('a.author', parent.author.name),
        m('div.text', parent.content),
        m('div.actions', [
          m('a.reply', {
            onclick: function() {
              module.vm.showForm(!module.vm.showForm());
            }
          }, 'Reply')
        ]),
        m('form.ui.reply.form.comment-form', {style: {display: module.vm.showForm() ? 'block' : 'none'}}, [
          m('textarea', {
            onchange: m.withAttr('value', module.vm.reply),
            value: module.vm.reply()
          }),
          m('button.ui.mini.blue.button', {onclick: module.vm.sendReply}, 'Reply')
        ])
      ]),
      m('div.comments', {
          style: {
            display: module.vm.childModules.length ? 'block' : 'none'
          }
        }, 
        module.vm.childModules.map(function(mod) {
          return mod.view();
        })
      )
    ])
  };

  return module;
};
