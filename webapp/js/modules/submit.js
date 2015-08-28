var app = app || {};

(function(module) {
  'use strict';

  var vm = {
    init: function() {
      this.errors = [];
      this.name = m.prop('');
    },
    submitTheme: function(e) {
      e.preventDefault();

      if (vm.errors.length) {
        vm.errors = [];
      }
      if (!vm.name().length) {
        vm.errors.push('no name');
        return
      }

      Theme.create({
        name: vm.name()
      }).then(
        function(theme) {
          m.route('/');
        }, 
        function(err) {
          if (err.status === 422) {
            vm.errors.push(err.msg);
          }
        });
    }
  };

  module.controller = function() {
    var user = storage.getUser();
    if (user === null) return m.route('/auth');

    vm.init();
  };

  module.view = function() {
    return [
      app.banner(),
      m('div.ui.grid', [
        app.usermenu(),
        m('div.ten.wide.column', [
          m('form.ui.form', [
            m('div.ui.red.message', {
              style: {
                display: vm.errors.length
                  ? 'block'
                  : 'none'
              }
            }, [
              m('ul.list', vm.errors.map(function(err) {
                return m('li', err);
              }))
            ]),
            m('div#tname.eight.wide.field', [
              m('input[type="text"][name="name"][placeholder="name"]', {
                onchange: m.withAttr('value', vm.name)
              })
            ]),
            m('div.field', [
              m('input.ui.tiny.blue.button[type="submit"][value="submit"]', {
                onclick: vm.submitTheme
              })
            ])
          ])
        ])
      ])
    ];
  };
})(app.newTheme = {});

(function(module) {
  'use strict';

  var LINK = 'link',
      TEXT = 'text';

  var vm = {
    init: function(t) {
      var user = storage.getUser();
      if (user === null) return m.route('/auth');

      this.errors = [];
      this.type = m.prop(t || LINK);
      this.title = m.prop('');
      this.link = m.prop('');
      this.content = m.prop('');
      this.theme = m.prop('');
      this.themes = m.prop([]);

      Theme.all().then(this.themes, app.utils.processError);
    },
    submitItem: function(e) {
      e.preventDefault();

      if (vm.errors.length) {
        vm.errors = [];
      }
      if (!vm.title().length) {
        vm.errors.push('no title');
      }
      if (!vm.theme().length) {
        vm.errors.push('choose a theme');
      }
      if (vm.errors.length) return;

      Item.create({
        type: vm.type(),
        title: vm.title(),
        link: vm.link(),
        content: vm.content(),
        theme: vm.theme()
      }).then(
      function(item) {
        console.log(item);
        m.route('/b/' + vm.theme());
      },
      function(err) {
        vm.errors.push(err.msg);
      });
    }
  }

  module.controller = function() {
    vm.init();
  }

  module.view = function(ctrl) {
    return [
      app.banner(),
      m('div.ui.grid', [
        app.usermenu(),
        m('div.ten.wide.column', [
          m('div.ui.tabular.menu', [
            m('a.item.active', {
              onclick: function() {
                vm.errors = [];
                vm.type(LINK);
                $('.active').removeClass('active');
                $(this).addClass('active');
              }
            }, 'Link '),
            m('a.item', {
              onclick: function() {
                vm.errors = [];
                vm.type(TEXT);
                $('.active').removeClass('active');
                $(this).addClass('active');
              }
            }, 'Text')
          ]),
          m('form.ui.form', {style: {width: '280px'}}, [
            m('div.ui.red.message', {
              style: {
                display: vm.errors.length ? 'block' : 'none'
              }
            }, [
              m('ul.list', vm.errors.map(function(err) {
                return m('li', err);
              }))
            ]),
            m('div.field', [
              m('input[type="text"][name="title"][placeholder="title"]', {
                onchange: m.withAttr('value', vm.title)
              })
            ]),
            m('div.field', [
              m('input[type="text"][name="link"][placeholder="link"]', {
                onchange: m.withAttr('value', vm.link),
                style: {display: vm.type() === LINK ? 'block' : 'none'}
              })
            ]),
            m('div.field', [
              m('textarea[name="content"]', {
                onchange: m.withAttr('value', vm.content),
                style: {display: vm.type() === TEXT ? 'block' : 'none'}
              })
            ]),
            m('div.field', [
              m('input[type="text"][name="brhub"][placeholder="theme"]', {
                value: vm.theme(),
                onchange: m.withAttr('value', vm.theme)
              })
            ]),
            m('p', vm.themes().map(function(theme) {
                return m('span', {
                  onclick: function() {
                    vm.theme(theme.name);
                  },
                  style: {
                    display: 'inline-block',
                    color: theme.color,
                    cursor: 'pointer',
                    'margin-right': '5px'
                  }
                }, theme.name);
              })
            ),
            m('input.ui.tiny.blue.button[type="submit"][value="submit"]', {
              onclick: vm.submitItem
            })
          ])
        ])
      ])
    ]
  }
})(app.newItem = {});
