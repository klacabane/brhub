var User = {};

User.signin = function(name, pw) {
  return Req({
    method: 'POST', 
    ep: '/auth', 
    data: {
      name: name,
      password: pw
    }
  }, true);
};

var Item = {};

Item.get = function(id) {
  return Req({
    method: 'GET',
    ep: '/api/items/' + id
  });
};

Item.create = function(data) {
  return Req({
    method: 'POST',
    ep: '/api/items/',
    data: data
  });
};

var Comments = {};
Comments.create = function(data) {
  console.log(data)
  return Req({
    method: 'POST',
    ep: '/api/comments/',
    data: data
  });
}

var Theme = {};

Theme.all = function() {
  return Req({
    method: 'GET',
    ep: '/api/b/'
  });
};

Theme.items = function(ep, start, n) {
  ep = ep === 'timeline' 
    ? 'timeline' 
    : 'b/' + ep;
  ep += '/' + start + '/' + n;

  return Req({
    method: 'GET',
    ep: '/api/' + ep
  });
};

Theme.create = function(data) {
  return Req({
    method: 'POST',
    ep: '/api/b/',
    data: data
  });
};


