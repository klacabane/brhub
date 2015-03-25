var app = app || {};

app.newBrhub = {
  controller: function() {
    var user = storage.getUser();
    if (user === null) return m.route('/auth');

    this.banner = app.banner(user);

    this.name = m.prop('');

    this.submitBrhub = function(e) {
      e.preventDefault();

      Brhub.create({
        name: this.name(),
      }).then(function(b) {
        console.log(b)
      }, 
      app.utils.processError);
    }.bind(this);
  },
  view: function(ctrl) {
    return [
      ctrl.banner.view(),
      m('div[class="content"]', [
        m('form[class="form-inline"]', [
          m('input[type="text"][class="form-control"][name="name"][placeholder="name"]', {
            onchange: m.withAttr('value', ctrl.name)
          }),
          m('input[type="submit"][value="submit"][class="btn btn-default"]', {onclick: ctrl.submitBrhub})
        ])
      ])
    ];
  }
};

(function(module) {
  'use strict';

  var LINK = 'link',
      TEXT = 'text';

  module.vm = {
    init: function(t) {
      var user = storage.getUser();
      if (user === null) return m.route('/auth');

      this.banner = app.banner(user);

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
  }

  module.controller = function() {
    module.vm.init();
  }

  module.view = function(ctrl) {
    return [
      module.vm.banner.view(),
      m('div[class="content"]', [
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
      ])
    ]
  }
})(app.submit = {});
