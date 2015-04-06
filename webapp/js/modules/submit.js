var app = app || {};

app.newTheme = {
  controller: function() {
    var user = storage.getUser();
    if (user === null) return m.route('/auth');

    this.name = m.prop('');
    this.submitTheme = function(e) {
      e.preventDefault();

      Theme.create({
        name: this.name()
      }).then(function(theme) {
        m.route('/');
      }, 
      app.utils.processError);
    }.bind(this);
  },
  view: function(ctrl) {
    return [
      app.banner(),
      m('div[class="ui grid"]', [
        app.usermenu(),
        m('div[class="ten wide column"]', [
          m('form[class="ui form"]', [
            m('div[class="eight wide field"]', [
              m('input[type="text"][name="name"][placeholder="name"]', {
                onchange: m.withAttr('value', ctrl.name)
              })
            ]),
            m('div[class="field"]', [
              m('input[type="submit"][value="submit"][class="ui tiny blue button"]', {onclick: ctrl.submitTheme})
            ])
          ])
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

      var vm = module.vm;
      Item.create({
        type: vm.type(),
        title: vm.title(),
        link: vm.link(),
        content: vm.content(),
        brhub: vm.theme()
      }).then(function(item) {
        console.log(item);
        m.route('/b/' + vm.theme());
      }, app.utils.processError);
    }
  }

  module.controller = function() {
    module.vm.init();
  }

  module.view = function(ctrl) {
    return [
      app.banner(),
      m('div[class="ui grid"]', [
        app.usermenu(),
        m('div[class="ten wide column"]', [
          m('div[class="ui tabular menu"]', [
            m('a[class="item active"]', {
              onclick: function() {
                module.vm.type(LINK);
                $('.active').removeClass('active');
                $(this).addClass('active');
              }
            }, 'Link '),
            m('a[class="item"]', {
              onclick: function() {
                module.vm.type(TEXT);
                $('.active').removeClass('active');
                $(this).addClass('active');
              }
            }, 'Text')
          ]),
          m('form[class="ui form"]', {style: {width: '280px'}}, [
            m('div[class="field"]', [
              m('input[type="text"][name="title"][placeholder="title"]', {
                onchange: m.withAttr('value', module.vm.title)
              })
            ]),
            m('div[class="field"]', [
              m('input[type="text"][name="link"][placeholder="link"]', {
                onchange: m.withAttr('value', module.vm.link),
                style: {display: module.vm.type() === LINK ? 'block' : 'none'}
              })
            ]),
            m('div[class="field"]', [
              m('textarea[name="content"]', {
                onchange: m.withAttr('value', module.vm.content),
                style: {display: module.vm.type() === TEXT ? 'block' : 'none'}
              })
            ]),
            m('div[class="field"]', [
              m('input[type="text"][name="brhub"][placeholder="theme"]', {
                value: module.vm.theme(),
                onchange: m.withAttr('value', module.vm.theme)
              })
            ]),
            m('p', module.vm.themes().map(function(theme) {
                return m('span', {
                  onclick: function() {
                    module.vm.theme(theme.name);
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
            m('input[type="submit"][class="ui tiny blue button"][value="submit"]', {onclick: module.vm.submitItem})
          ])
        ]),
      ])
    ]
  }
})(app.newItem = {});
