var app = app || {};

(function(module) {
  'use strict';

  var sort = function(list) {
    return {
      onclick: function(e) {
        var prop = e.target.getAttribute('data-sort-by'),
        first = list[0],
        fn;
      
        if (prop) {
          var props = prop.split('.');
          if (props.length === 2) {
            fn = function(a, b) { 
              return a[props[0]][props[1]] > b[props[0]][props[1]] ? 1 : a[props[0]][props[1]] < b[props[0]][props[1]] ? -1 : 0; 
            }
          } else {
            fn = function(a, b) {
              return a[props[0]] > b[props[0]] ? 1 : a[props[0]] < b[props[0]] ? -1 : 0;
            }
          }
          list.sort(fn);
          if (first === list[0]) list.reverse();
        }
      }
    }
  }

  module.vm = {
    src: '',
    current: 0,
    limit: 5,
    showPrev: false,
    showNext: false,
    items: m.prop([]),
    getItems: function() {
      Brhub.items(this.src, this.current, this.limit)
        .then(function(res) {
          this.showPrev = this.current > 0;
          this.showNext = res.hasmore;
          this.items(res.items);
        }.bind(this), app.utils.processError);
    },
    getPrevItems: function() {
      this.current -= this.limit;
      if (this.current < 0) this.current = 0;

      this.getItems();
    },
    getNextItems: function() {
      this.current += this.limit;
      this.getItems();
    }
  };

  module.controller = function(opts) {
    module.vm.src = opts.src;
    module.vm.current = opts.start || 0;

    module.vm.getItems();
  }

  module.view = function(ctrl) {
    return m('div', [
      m('table', [
        m('tr', sort(module.vm.items()), [
          m('th[data-sort-by=name]', 'Title'),
          m('th[data-sort-by=date]', 'Date'),
          m('th[data-sort-by=brhub.name]', 'Brhub')
        ]),
        module.vm.items().map(function(item) {
          return m("tr", {
            onclick: function(e) {
              // m.route('/items/' + item.id)
            }}, [
            m('td', item.title),
            m('td', item.date),
            m('td', item.brhub.name)
          ]) 
        })
      ]),
      m('button', {
        onclick: module.vm.getPrevItems.bind(module.vm), 
        style: {
          display: module.vm.showPrev ? 'block' : 'none'
        }}, 'Prev'),
      m('button', {
        onclick: module.vm.getNextItems.bind(module.vm), 
        style: {
          display: module.vm.showNext ? 'block' : 'none'
        }}, 'Next')
    ]);
  }
})(app.grid = {})
