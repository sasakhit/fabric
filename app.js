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
var helper = require('./utils/helper.js');
var channels = require('./utils/create-channel.js');
var join = require('./utils/join-channel.js');
var install = require('./utils/install-chaincode.js');
var instantiate = require('./utils/instantiate-chaincode.js');
var invoke = require('./utils/invoke-transaction.js');
var query = require('./utils/query.js');
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var util = require('util');
var fs = require('fs');
const https = require('https');

var host = process.env.HOST || config.host;
var port = process.env.PORT || config.port;
var peers = config.peers;
logger.debug('peers : ' + config.peers);

// Route includes
//var login = require('./routes/login');
//var helloworld = require('./routes/helloworld');
//var counter = require('./routes/counter');
//var fx = require('./routes/fx');
//var sl = require('./routes/sl');

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

// Serve back static files
app.use(express.static('public'));
app.use(express.static('public/views'));

// App Set //
app.set('port', (port || 3000));

app.use(function(req, res, next) {
  req.username = app.get('username');
  req.orgname = app.get('orgname');

  return next();
});

// Listen //
console.log("Staring http server on: " + app.get("port"));
var server = http.createServer(app).listen(app.get("port"), function () {
    console.log("Server Up - " + app.get("port"));
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

// HFC
//var hfc = require('hfc');

///////////////////////////////////////////////////////////////////////////////
///////////////////////// REST ENDPOINTS START HERE ///////////////////////////
///////////////////////////////////////////////////////////////////////////////
// Register and enroll user
app.post('/users', function(req, res) {
	var username = req.body.username;
	var orgname = req.body.orgname;
  app.set('username', username);
  app.set('orgname', orgname);

	logger.debug('End point : /users');
	logger.debug('User name : ' + username);
	logger.debug('Org name  : ' + orgname);
	if (!username) {
		res.json(getErrorMessage('\'username\''));
		return;
	}
	if (!orgname) {
		res.json(getErrorMessage('\'orgname\''));
		return;
	}

	helper.getRegisteredUsers(username, orgname, true).then(function(response) {
		if (response && typeof response !== 'string') {
			//response.token = token;
			res.json(response);
		} else {
			res.json({
				success: false,
				message: response
			});
		}
	});

});
// Create Channel
app.post('/channels', function(req, res) {
	logger.info('<<<<<<<<<<<<<<<<< C R E A T E  C H A N N E L >>>>>>>>>>>>>>>>>');
	logger.debug('End point : /channels');
	var channelName = req.body.channelName;
	var channelConfigPath = req.body.channelConfigPath;
	logger.debug('Channel name : ' + channelName);
	logger.debug('channelConfigPath : ' + channelConfigPath); //../artifacts/channel/mychannel.tx
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!channelConfigPath) {
		res.json(getErrorMessage('\'channelConfigPath\''));
		return;
	}

	channels.createChannel(channelName, channelConfigPath, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Join Channel
app.post('/channels/:channelName/peers', function(req, res) {
	logger.info('<<<<<<<<<<<<<<<<< J O I N  C H A N N E L >>>>>>>>>>>>>>>>>');
	var channelName = req.params.channelName;
	//var peers = req.body.peers;
	logger.debug('channelName : ' + channelName);
	logger.debug('peers : ' + peers);
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
		return;
	}

	join.joinChannel(channelName, peers, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Install chaincode on target peers
app.post('/chaincodes', function(req, res) {
	logger.debug('==================== INSTALL CHAINCODE ==================');
	//var peers = req.body.peers;
	var chaincodeName = req.body.chaincodeName;
	var chaincodePath = req.body.chaincodePath;
	var chaincodeVersion = req.body.chaincodeVersion;
	logger.debug('peers : ' + peers); // target peers list
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('chaincodePath  : ' + chaincodePath);
	logger.debug('chaincodeVersion  : ' + chaincodeVersion);
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
		return;
	}
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!chaincodePath) {
		res.json(getErrorMessage('\'chaincodePath\''));
		return;
	}
	if (!chaincodeVersion) {
		res.json(getErrorMessage('\'chaincodeVersion\''));
		return;
	}

	install.installChaincode(peers, chaincodeName, chaincodePath, chaincodeVersion, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Instantiate chaincode on target peers
app.post('/channels/:channelName/chaincodes', function(req, res) {
	logger.debug('==================== INSTANTIATE CHAINCODE ==================');
	var chaincodeName = req.body.chaincodeName;
	var chaincodeVersion = req.body.chaincodeVersion;
	var channelName = req.params.channelName;
	var functionName = req.body.functionName;
	var args = req.body.args;
  peers = config.peers_instantiate;
  logger.debug('peers  : ' + peers);
	logger.debug('channelName  : ' + channelName);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('chaincodeVersion  : ' + chaincodeVersion);
	logger.debug('functionName  : ' + functionName);
	logger.debug('args  : ' + args);
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!chaincodeVersion) {
		res.json(getErrorMessage('\'chaincodeVersion\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!functionName) {
		res.json(getErrorMessage('\'functionName\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}
	instantiate.instantiateChaincode(channelName, chaincodeName, chaincodeVersion, functionName, args, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Invoke transaction on chaincode on target peers
app.post('/channels/:channelName/chaincodes/:chaincodeName', function(req, res) {
	logger.debug('==================== INVOKE ON CHAINCODE ==================');
	//var peers = req.body.peers;
	var chaincodeName = req.params.chaincodeName;
	var channelName = req.params.channelName;
	var fcn = req.body.fcn;
	var args = req.body.args;
  peers = config.peers;
	logger.debug('channelName  : ' + channelName);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
		return;
	}
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!fcn) {
		res.json(getErrorMessage('\'fcn\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}

	invoke.invokeChaincode(peers, channelName, chaincodeName, fcn, args, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Query on chaincode on target peers
app.get('/channels/:channelName/chaincodes/:chaincodeName', function(req, res) {
	logger.debug('==================== QUERY BY CHAINCODE ==================');
	var channelName = req.params.channelName;
	var chaincodeName = req.params.chaincodeName;
	let args = req.query.args;
	let fcn = req.query.fcn;
	let peer = req.query.peer;

	logger.debug('channelName : ' + channelName);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('fcn : ' + fcn);
	logger.debug('args : ' + args);

	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!fcn) {
		res.json(getErrorMessage('\'fcn\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}
	args = args.replace(/'/g, '"');
	args = JSON.parse(args);
	logger.debug(args);

	query.queryChaincode(peer, channelName, chaincodeName, args, fcn, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
//  Query Get Block by BlockNumber
app.get('/channels/:channelName/blocks/:blockId', function(req, res) {
	logger.debug('==================== GET BLOCK BY NUMBER ==================');
	let blockId = req.params.blockId;
	let peer = req.query.peer;
	logger.debug('channelName : ' + req.params.channelName);
	logger.debug('BlockID : ' + blockId);
	logger.debug('Peer : ' + peer);
	if (!blockId) {
		res.json(getErrorMessage('\'blockId\''));
		return;
	}

	query.getBlockByNumber(peer, blockId, req.username, req.orgname)
		.then(function(message) {
			res.send(message);
		});
});
// Query Get Transaction by Transaction ID
app.get('/channels/:channelName/transactions/:trxnId', function(req, res) {
	logger.debug(
		'================ GET TRANSACTION BY TRANSACTION_ID ======================'
	);
	logger.debug('channelName : ' + req.params.channelName);
	let trxnId = req.params.trxnId;
	let peer = req.query.peer;
	if (!trxnId) {
		res.json(getErrorMessage('\'trxnId\''));
		return;
	}

	query.getTransactionByID(peer, trxnId, req.username, req.orgname)
		.then(function(message) {
			res.send(message);
		});
});
// Query Get Block by Hash
app.get('/channels/:channelName/blocks', function(req, res) {
	logger.debug('================ GET BLOCK BY HASH ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let hash = req.query.hash;
	let peer = req.query.peer;
	if (!hash) {
		res.json(getErrorMessage('\'hash\''));
		return;
	}

	query.getBlockByHash(peer, hash, req.username, req.orgname).then(
		function(message) {
			res.send(message);
		});
});
//Query for Channel Information
app.get('/channels/:channelName', function(req, res) {
	logger.debug(
		'================ GET CHANNEL INFORMATION ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let peer = req.query.peer;

  if (req.username && req.orgname) {
	  query.getChainInfo(peer, req.username, req.orgname).then(
		  function(message) {
			  res.send(message);
		  });
  }
});
// Query to fetch all Installed/instantiated chaincodes
app.get('/chaincodes', function(req, res) {
	var peer = req.query.peer;
	var installType = req.query.type;
	//TODO: add Constnats
	if (installType === 'installed') {
		logger.debug(
			'================ GET INSTALLED CHAINCODES ======================');
	} else {
		logger.debug(
			'================ GET INSTANTIATED CHAINCODES ======================');
	}

	query.getInstalledChaincodes(peer, installType, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Query to fetch channels
app.get('/channels', function(req, res) {
	logger.debug('================ GET CHANNELS ======================');
	logger.debug('peer: ' + req.query.peer);
	var peer = req.query.peer;
	if (!peer) {
		res.json(getErrorMessage('\'peer\''));
		return;
	}

	query.getChannels(peer, req.username, req.orgname)
	.then(function(
		message) {
		res.send(message);
	});
});

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
        console.error('chainstats failure: (' +
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

/*
function init() {
	if (isHSBN) {
		certFile = '0.secure.blockchain.ibm.com.cert';
	}
	fs.createReadStream(certFile).pipe(fs.createWriteStream(certPath));
	enrollAndRegisterUsers();
}
*/

/*
function enrollAndRegisterUsers() {
    var cert = fs.readFileSync(certFile);

    chain.setMemberServicesUrl(ca_url, {
        pem: cert
    });

    // Adding all the peers to blockchain
    // this adds high availability for the client
    for (var i = 0; i < peers.length; i++) {

        // Peers on Bluemix require secured connections, hence 'grpcs://'
        chain.addPeer("grpcs://" + peers[i].discovery_host + ":" + peers[i].discovery_port, {
            pem: cert
        });
    }

    console.log("\n\n------------- peers and caserver information: -------------");
    console.log(chain.getPeers());
    console.log(chain.getMemberServices());
    console.log('-----------------------------------------------------------\n\n');

    //setting timers for fabric waits
    chain.setDeployWaitTime(config.deployWaitTime);
    app.set('chain', chain);

    // Get chain information
    var url = 'https://' + peers[0].api_host + ':' + peers[0].api_port + '/chain';
    https.get(url, function(res){
      var body = '';
      res.setEncoding('utf8');

      res.on('data', function(chunk){
        body += chunk;
      });

      res.on('end', function(res){
        ret = JSON.parse(body);
        console.log(ret);
      });
    }).on('error', function(e){
      console.log(e.message);
    });

    // Monitor chain's blockheight and pass it along to clients.
    setInterval(function () {
      var options = {
        host: peers[0].api_host,
        port: peers[0].api_port,
        path: '/chain',
        method: 'GET'
      };

      function success(statusCode, headers, resp) {
        resp = JSON.parse(resp);
        if (resp && resp.height) {
          io.sockets.emit('block_change', resp);
        }
      }

      function failure(statusCode, headers, msg) {
      // Don't broadcast failures to clients, just log them
        console.error('chainstats failure: (' +
                      'status code: ' + statusCode +
                      '\n  headers: ' + headers +
                      '\n  message: ' + msg + ')');
      }

      var request = https.request(options, function (resp) {
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
*/
