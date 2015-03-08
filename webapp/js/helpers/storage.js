var storage = {
  lsKey: 'brhub-user',
  getUser: function() {
    var data = localStorage.getItem(this.lsKey);
    var user = null;
    if (data === null) {
      return null;
    }

    try {
      user = JSON.parse(localStorage.getItem(this.lsKey));
    } catch (e) {
      localStorage.removeItem(this.lsKey);
    } 
    return user;
  },
  setUser: function (user) {
    if (user === null)
      localStorage.removeItem(this.lsKey);
    else
      localStorage.setItem(this.lsKey, JSON.stringify(user));
  }
}
