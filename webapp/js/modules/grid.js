var app = app || {};

(function(module) {
  'use strict';

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
      module.vm.items().map(function(item) {
        return m("div", [
          m('p',[
            m('a', {
              href: item.type === 'link' ? item.link : '/items/' + item.id
            }, item.title)
          ]),
          m('p', [
            m('span', 'by ' + item.author.name),
            m('span', {
              onclick: function(e) {
                m.route('/b/' + item.brhub)
              },
              style: {display: module.vm.src === 'timeline' ? 'inline' : 'none'}
            }, item.brhub)
          ])
        ]) 
      }),
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

  var sort = function(list) {
    return {
      onclick: function(e) {
        var prop = e.target.getAttribute('data-sort-by'),
        first = list[0];
      
        if (prop) {
          list.sort(function(a, b) {
            return a[prop] > b[prop] ? 1 : a[prop] < b[prop] ? -1 : 0;
          });
          if (first === list[0]) list.reverse();
        }
      }
    }
  }
})(app.grid = {});
