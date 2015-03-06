var app = app || {};

(function (module) {
  'use strict';

  var sort = function (list) {
    return {
      onclick: function (e) {
        var prop = e.target.getAttribute('data-sort-by'),
        first = list[0],
        fn;
      
        if (prop) {
          var props = prop.split('.');
          if (props.length === 2) {
            fn = function (a, b) { 
              return a[props[0]][props[1]] > b[props[0]][props[1]] ? 1 : a[props[0]][props[1]] < b[props[0]][props[1]] ? -1 : 0; 
            }
          } else {
            fn = function (a, b) {
              return a[props[0]] > b[props[0]] ? 1 : a[props[0]] < b[props[0]] ? -1 : 0;
            }
          }
          list.sort(fn);
          if (first === list[0]) list.reverse();
        }
      }
    }
  }

  module.controller = function () {
    var ctrl = this;
    this.items = m.prop();

    User.timeline().then(function (data) {
      data.sort(function (a, b) {
        return a.date > b.date ? 1 : a.date < b.date ? -1 : 0;
      });
      ctrl.items(data);
    }, app.utils.processError);
  }

  module.view = function (ctrl) {
    return m('table', [
      m('tr', sort(ctrl.items()), [
        m('th[data-sort-by=date]', 'Date'),
        m('th[data-sort-by=brhub.name]', 'Brhub')
      ]),
      ctrl.items().map(function (item) {
        return m("tr", {
          onclick: function (e) {
            m.route('/items/' + item.id)
          }}, [
          m('td', item.id),
          m('td', item.name),
          m('td', item.brhub.name)
        ]) 
      })
    ]);
  }
})(app.grid = {})
