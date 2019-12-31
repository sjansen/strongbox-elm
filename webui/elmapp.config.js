const proxy = require('http-proxy-middleware');

module.exports = {
  setupProxy: function(app) {
    app.use(proxy('/api', { target: 'http://localhost:8080/' }));
  }
};
