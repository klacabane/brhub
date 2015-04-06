var app = app || {};

var commentModule = function(parent) {
  var module = {};

  module.vm = {
    reply: m.prop(''),
    showForm: m.prop(false),
    sendReply: function(e) {
      e.preventDefault();

      if (!module.vm.reply().length)
        return;

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
    return m('div[class="comment"]', [
      m('a[class="avatar"]', [
        m('img', {src: 'images/avatar.jpg'})
      ]),
      m('div[class="content"]', [
        m('a[class="author"]', parent.author.name),
        m('div[class="text"]', parent.content),
        m('div[class="actions"]', [
          m('a[class="reply"]', {
            onclick: function() {
              module.vm.showForm(!module.vm.showForm());
            }
          }, 'Reply')
        ]),
        m('form[class="ui reply form comment-form"]', {style: {display: module.vm.showForm() ? 'block' : 'none'}}, [
          m('textarea', {
            onchange: m.withAttr('value', module.vm.reply),
            value: module.vm.reply()
          }),
          m('button[class="ui mini blue button"]', {onclick: module.vm.sendReply}, 'Reply')
        ])
      ]),
      m('div[class="comments"]', {
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
