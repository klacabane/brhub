(function(app) {
  'use strict';

  app.utils = {
    processError: function(err) {
      console.log(err)
      switch (err.status) {
        case 401:
        case 403:
          storage.setUser(null);
          m.route('/auth');
          break;
        case 500:
      }
    }
  }

  m.route.mode = 'hash';
  m.route(document.body, '/', {
    '/auth': app.auth,
    '/': app.feed('timeline'),
    '/b/:name': app.feed(),
    '/items/:id': app.items,
    '/submit': app.newItem,
    '/theme': app.newTheme,
    '/users/:name': app.profile
  });
})(app = app || {});
