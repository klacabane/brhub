(function(app) {
  'use strict';

  app.utils = {
    processError: function(err) {
      console.log(err)
      switch (err.status) {
        case 403:
          storage.setUser(null);
          m.route('/');
          break;
        case 500:
      }
    }
  }

  m.route.mode = 'hash';
  m.route(document.body, '/', {
    '/': app.auth,
    '/timeline': app.timeline
  });
})(app = app || {});
