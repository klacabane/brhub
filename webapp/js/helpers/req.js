var Req = function(opts, preventAuth)  {
  return m.request({
    method: opts.method,
    url: 'http://localhost:8000' + opts.ep,
    data: opts.data || {},
    config: function(xhr) {
      xhr.setRequestHeader('Content-Type', 'application/json');
    
      if (!preventAuth) {
        var tkn = storage.getUser().token;
        xhr.setRequestHeader('X-token', tkn);
      }
    },
    extract: function(xhr) {
      return xhr.status >= 400 
        ? JSON.stringify({err: JSON.parse(xhr.responseText).msg, status: xhr.status})
        : xhr.responseText;
    }
  })
}
