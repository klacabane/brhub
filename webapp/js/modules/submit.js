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
    submitItem: function(e) {
      e.preventDefault();

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
        m('p[class="text-primary"]', [
          m('span', {
            onclick: function() {
              module.vm.type(LINK);
            }
          }, 'Link '),
          m('span', {
            onclick: function() {
              module.vm.type(TEXT);
            }
          }, 'Text')
        ]),
        m('form', {style: {width: '280px'}}, [
          m('input[type="text"][class="form-control"][name="title"][placeholder="title"]', {
            onchange: m.withAttr('value', module.vm.title)
          }),
          m('input[type="text"][class="form-control"][name="link"][placeholder="link"]', {
            onchange: m.withAttr('value', module.vm.link),
            style: {display: module.vm.type() === LINK ? 'block' : 'none'}
          }),
          m('textarea[class="form-control"][name="content"]', {
            onchange: m.withAttr('value', module.vm.content),
            style: {display: module.vm.type() === TEXT ? 'block' : 'none'}
          }),
          m('input[type="text"][class="form-control"][name="brhub"][placeholder="theme"]', {
            value: module.vm.brhub(),
            onchange: m.withAttr('value', module.vm.brhub)
          }),
          m('p', module.vm.brhubs().map(function(b) {
              return m('span', {
                onclick: function() {
                  module.vm.brhub(b.name);
                },
                style: {
                  display: 'inline-block',
                  color: b.color,
                  cursor: 'pointer',
                  'margin-right': '5px'
                }
              }, b.name);
            })
          ),
          m('input[type="submit"][class="btn btn-default"][value="submit"]', {onclick: module.vm.submitItem})
        ]),
      ])
    ]
  }
})(app.submit = {});
