(function(app) {
  'use strict';

  app.utils = {
    processError: function (err) {
      console.log("ERROR:", err.msg);
    }
  }

  m.route.mode = 'hash';
  m.route(document.body, '/', {
    '/': app.home,
    '/dashboard': app.dashboard
  });
})(app = app || {});
