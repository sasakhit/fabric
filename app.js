process.env.GOPATH = __dirname

// Web
var express = require('express');
var app = express();
var path = require('path');
var bodyParser = require('body-parser');
var session = require('express-session');
var http = require('http');
var socketio = require('socket.io');
var winston = require('winston');								//logginer module
var os = require('os');
var expressJWT = require('express-jwt');
var jwt = require('jsonwebtoken');
var bearerToken = require('express-bearer-token');
var cors = require('cors');
var config = require('./config.json');
var log4js = require('log4js');
var logger = log4js.getLogger('Main');
var util = require('util');
var fs = require('fs');
const https = require('https');

var host = process.env.HOST || config.host;
var port = process.env.PORT || config.port;

// Route includes
//var login = require('./routes/login');
//var helloworld = require('./routes/helloworld');
//var counter = require('./routes/counter');
//var fx = require('./routes/fx');
//var sl = require('./routes/sl');
var fabric = require('./routes/fabric');

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));

// Passport Session Configuration //
app.use(express.static(path.join(__dirname, 'public')));
app.use(session({
    secret: 'secret',
    key: 'user',
    resave: 'true',
    saveUninitialized: false,
    cookie: {maxage: 60000, secure: false}
}));

//app.use('/login', login);
//app.use('/helloworld', helloworld);
//app.use('/counter', counter);
//app.use('/fx', fx);
//app.use('/sl', sl);
app.use('/', fabric);

// Serve back static files
app.use(express.static('public'));
app.use(express.static('public/views'));

// App Set //
app.set('port', (port || 3000));

// Listen //
logger.debug("Staring http server on: " + app.get("port"));
var server = http.createServer(app).listen(app.get("port"), function () {
    logger.debug("Server Up - " + app.get("port"));
});

function getErrorMessage(field) {
	var response = {
		success: false,
		message: field + ' field is missing or Invalid in the request'
	};
	return response;
}

// Link socket to http server //
var io = socketio.listen(server);
//monitorBlockHeight();

// Load socket module for chat
require('./routes/chat')(io);
require('./routes/vote')(io);


function monitorBlockHeight() {

    setInterval(function () {

      var options = {
        host: 'localhost',
        port: '4000',
        path: '/channels/mychannel',
        params: {peer: 'peer1'},
        method: 'GET'
      };

      function success(statusCode, headers, resp) {
        resp = JSON.parse(resp);
        if (resp && resp.height.low) {
          io.sockets.emit('block_change', resp);
        }
      }

      function failure(statusCode, headers, msg) {
      // Don't broadcast failures to clients, just log them
        logger.error('chainstats failure: (' +
                     'status code: ' + statusCode +
                     '\n  headers: ' + headers +
                     '\n  message: ' + msg + ')');
      }

      var request = http.request(options, function (resp) {
        var str = '', chunks = 0;
        resp.setEncoding('utf8');

        resp.on('data', function (chunk) {                                                            //merge chunks of request
          str += chunk;
          chunks++;
        });
        resp.on('end', function () {
          if (resp.statusCode == 204 || resp.statusCode >= 200 && resp.statusCode <= 399) {
            success(resp.statusCode, resp.headers, str);
          }
          else {
            failure(resp.statusCode, resp.headers, str);
          }
        });
      });

      request.on('error', function (e) {
        failure(500, null, e);
      });

      request.setTimeout(20000);
      request.on('timeout', function () {                                                                //handle time out event
        failure(408, null, 'Request timed out');
      });

      request.end();
    }, 5000);
}
