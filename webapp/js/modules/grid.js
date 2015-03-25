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
  }

  module.controller = function(opts) {
    module.vm.src = opts.src;
    module.vm.current = opts.start || 0;

    module.vm.getItems();
  }

  module.view = function(ctrl) {
    return m('div[class="item-container"]', [
      module.vm.items().map(function(item) {
        return m('div[class="item-row"]', [
          m('a[class="item-image"]', {
            href: item.type === 'link' ? item.link : '/#/items/'+item.id
          }, [
            m('img', {
              src: '',
              width: 50,
              height: 50
            })
          ]),
          m('div[class="item-meta"]', [
            m('a[style="display: block; font-size: 16px; margin-bottom: 5px;"]', {href: item.type === 'link' ? item.link : '/#/items/'+item.id}, item.title),
            m('p', [
              m('em', 'by '+item.author.name+' '),
              m('strong', {
                onclick: function(e) {
                  m.route('/b/'+item.brhub.name);
                },
                style: {
                  display: module.vm.src === 'timeline' ? 'inline' : 'none',
                  color: item.brhub.color,
                  cursor: 'pointer',
                }
              }, item.brhub.name)
            ])
          ]),
          m('p', [
            m('a[style="font-size: 13px;"]', {
              href: '/#/items/'+item.id,
            }, item.commentCount === 1 ? item.commentCount + ' comment' : item.commentCount + ' comments')
          ])
        ])
      }),
      m('button', {
        onclick: module.vm.getPrevItems.bind(module.vm), 
        style: {
          display: module.vm.showPrev ? 'inline' : 'none'
        }}, 'Prev'),
      m('button', {
        onclick: module.vm.getNextItems.bind(module.vm), 
        style: {
          display: module.vm.showNext ? 'inline' : 'none'
        }}, 'Next')
    ]);
  }
})(app.grid = {});
