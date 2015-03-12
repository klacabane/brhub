var app = app || {};

(function(module) {
  'use strict';

  var LINK = 'link',
      TEXT = 'text';

  module.vm = {
    init: function(t) {
      if (storage.getUser() === null) return m.route('/auth');

      this.type = m.prop(t || LINK);
      this.title = m.prop('');
      this.link = m.prop('');
      this.content = m.prop('');
      this.brhub = m.prop('');
      this.brhubs = m.prop([]);

      Brhub.all().then(this.brhubs, app.utils.processError);
    },
    submitItem: function() {
      var vm = module.vm;

      Item.create({
        type: vm.type(),
        title: vm.title(),
        link: vm.link(),
        content: vm.content(),
        brhub: vm.brhub()
      }).then(function(item) {
        console.log(item);
        m.route('/b/' + vm.brhub());
      }, app.utils.processError);
    }
  };

  module.controller = function() {
    module.vm.init();
  };

  module.view = function(ctrl) {
    return m('div', [
      m('ul', [
        m('li', {
          onclick: function() {
            module.vm.type(LINK);
          }
        }, 'Link'),
        m('li', {
          onclick: function() {
            module.vm.type(TEXT);
          }
        }, 'Text')
      ]),
      m('input[name="title"][type="text"]', {
        onchange: m.withAttr('value', module.vm.title)
      }),
      m('input[name="link"][type="text"]', {
        onchange: m.withAttr('value', module.vm.link),
        style: {display: module.vm.type() === LINK ? 'block' : 'none'}
      }),
      m('textarea[name="content"]', {
        onchange: m.withAttr('value', module.vm.content),
        style: {display: module.vm.type() === TEXT ? 'block' : 'none'}
      }),
      m('input[name="brhub"][type="text"]', {value: module.vm.brhub()}),
      module.vm.brhubs().map(function(b) {
        return m('span', {
          onclick: function() {
            module.vm.brhub(b.name);
          }
        }, b.name);
      }),
      m('input[type="submit"][value="submit"]', {onclick: module.vm.submitItem})
    ]);
  };
})(app.submit = {});
