var User = {};

User.signin = function (name, pw) {
  return Req({
    method: 'POST', 
    ep: '/auth', 
    data: {
      name: name,
      password: pw
    }
  }, true);
};

User.timeline = function () {
  return function () {
    return Req({method: 'GET', ep: '/api/timeline/0/5'});
  }
}
