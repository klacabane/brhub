var app = app || {};

var commentModule = function(parent, root) {
  var module = {};
  var isItem = parent.hasOwnProperty('brhub');

  module.vm = {
    margin: (isItem || root) ? 0 : 40,
    reply: m.prop(''),
    showForm: m.prop(isItem),
    sendReply: function(e) {
      e.preventDefault();

      Comments.create({
        item: isItem ? parent.id : parent.item,
        parent: isItem ? '' : parent.id,
        content: module.vm.reply()
      }).then(function(comment) {
        module.vm.childModules.splice(0, 0, new commentModule(comment, isItem));
        module.vm.reply('');
      }, app.utils.processError);
    },
    childModules: parent.comments.map(function(child) {
      return new commentModule(child, isItem);
    })
  };

  module.view = function() {
    var childs = [];
    if (!isItem)
      childs.push(
        m('p', parent.content),
        m('ul', [
          m('li', {
            onclick: function() {
              module.vm.showForm(!module.vm.showForm());
            }
          }, 'Reply')
        ])
      );
 
    childs.push(
      m('form', {style: {display: module.vm.showForm() ? 'block' : 'none'}}, [
        m('textarea[name="content"]', {
          onchange: m.withAttr('value', module.vm.reply),
          value: module.vm.reply(),
        }),
        m('input[type="submit"][value="reply"]', {
          onclick: module.vm.sendReply,
        })
      ]),
      m('div', module.vm.childModules.map(function(sub) {
        return sub.view();
      }))
    );

    return m('div', {style: {marginLeft: module.vm.margin+'px'}}, childs);
  };

  return module;
};
