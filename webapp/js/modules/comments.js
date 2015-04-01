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

      if (!module.vm.reply().length)
        return;

      Comments.create({
        item: isItem ? parent.id : parent.item,
        parent: isItem ? '' : parent.id,
        content: module.vm.reply()
      }).then(function(comment) {
        module.vm.childModules.splice(0, 0, new commentModule(comment, isItem));
        module.vm.reply('');
        if (!isItem) {
          module.vm.showForm(false);
        }
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
        m('a', {href: '/#/users/'+parent.author.name}, parent.author.name),
        m('p[class="comment-content"]', parent.content),
        m('div', [
          m('small[class="btn-link"][style="cursor: pointer;"]', {
            onclick: function() {
              module.vm.showForm(!module.vm.showForm());
            }
          }, 'Reply')
        ])
      );
 
    childs.push(
      m('form', {style: {display: module.vm.showForm() ? 'block' : 'none'}}, [
        m('div[class="form-group"]', [
          m('textarea[style="width: 200px; margin-bottom: 2px;"][class="form-control"][name="content"]', {
            onchange: m.withAttr('value', module.vm.reply),
            value: module.vm.reply()
          }),
          m('button[type="submit"][class="btn btn-default btn-xs"]', {
            onclick: module.vm.sendReply
          }, 'Reply')
        ])
      ]),
      m('div', module.vm.childModules.map(function(sub) {
        return sub.view();
      }))
    );

    return m('div[class="comment"]', {style: {marginLeft: module.vm.margin+'px'}}, childs);
  };

  return module;
};
