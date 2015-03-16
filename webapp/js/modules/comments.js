var app = app || {};

var comments = function(parent, root) {
  var module = {};
  var isItem = parent.hasOwnProperty('brhub');
  module.vm = {
    margin: (isItem || root) ? 0 : 40,
    init: function() {
      var opts;
      if (isItem) {
        opts = {item: parent.id};
      } else {
        opts = {parent: parent.id, item: parent.item};
      }

      this.reply = new reply(opts);
      this.childModules = parent.comments.map(function(child) {
        return new comments(child, isItem);
      });
    }
  };

  module.view = function() {
    var childs = [];
    if (!isItem)
      childs.push(m('p', parent.content));
 
    childs.push(
      m('ul', [
        m('li', 'Reply')
      ]),
      module.vm.reply.view(),
      m('div', module.vm.childModules.map(function(sub) {
        return sub.view();
      }))
    );

    return m('div', {style: {marginLeft: module.vm.margin+'px'}}, childs);
  };

  module.vm.init();
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
        item: module.vm.item,
        parent: module.vm.parent,
        content: module.vm.content()
      }).then(function(comment) {
        module.vm.content('');
        console.log(comment);
      }, app.utils.processError);
    }
  };

  module.view = function() {
    return m('form', [
      m('textarea[name="content"]', {
        onchange: m.withAttr('value', module.vm.content),
        value: module.vm.content(),
      }),
      m('input[type="submit"][value="reply"]', {
        onclick: module.vm.send,
      })
    ]);
  };
  return module;
};
