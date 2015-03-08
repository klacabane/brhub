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

User.newB = function() {
  return Req({
    method: 'POST',
    ep: '/api/b/',
    data: {
      name: 'f'
    }
  });
};

User.newItem = function(b) {
  return Req({
    method: 'POST',
    ep: '/api/items/',
    data: {
      content: 'axaxa',
      brhub: "54faf021d43c6c040b000001"
    }
  });
}

var Brhub = {};

Brhub.items = function(ep, start, n) {
  ep = ep === 'timeline' 
    ? 'timeline' 
    : 'b/' + ep;
  ep += '/' + start + '/' + n;

  return Req({
    method: 'GET',
    ep: '/api/' + ep
  });
}
