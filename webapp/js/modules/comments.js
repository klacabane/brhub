var app = app || {};

var comments = function(parent, root) {
  var module = {};
  module.vm = {
    reply: new reply({parent: parent.id, item: parent.item}),
    parent: parent,
    margin: root ? 0 : 40,
    childs: parent.comments.map(function(child) {
      return new comments(child);
    })
  };

  module.view = function() {
    return m('div', {style: {marginLeft: module.vm.margin+'px'}}, [
      m('p', module.vm.parent.content),
      m('ul', [
        m('li', 'Reply')
      ]),
      module.vm.reply.view(),
      m('div', module.vm.childs.map(function(sub) {
        return sub.view();
      }))
    ]);
  };
  return module;
};

var reply = function(opts) {
  var module = {};
  module.vm = {
    parent: opts.parent || '',
    item: opts.item,
    content: m.prop(''),
    send: function(e) {
      e.preventDefault();

      Comments.create({
        item: this.item,
        parent: this.parent,
        content: this.content()
      }).then(function(comment) {
        console.log(comment);
      }, app.utils.processError);
    }
  };

  module.view = function() {
    return m('form', [
      m('textarea[name="content"]', {
        onchange: m.withAttr('value', module.vm.content)
      }),
      m('input[type="submit"][value="reply"]', {
        onclick: module.vm.send.bind(module.vm)
      })
    ]);
  };
  return module;
};
