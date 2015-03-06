var Req = function (opts, preventAuth)  {
  return m.request({
    method: opts.method,
    url: 'http://localhost:8000' + opts.ep,
    data: opts.data || {},
    config: function (xhr) {
      xhr.setRequestHeader('Content-Type', 'application/json');
    
      if (!preventAuth) {
        var tkn = storage.getUser().token;
        xhr.setRequestHeader('X-token', tkn);
      }
    }
  })
}
