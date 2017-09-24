var express = require('express');
var router = express.Router();

router.get('/', function(req, res) {
  var username = req.query.username;
  var password = req.query.password;
  var chain = req.app.get('chain');
  chain.enroll(username, password, function (err, user) {
    if (err) return res.json("ERROR: failed to enroll user : " + err);

    console.log("\nEnrolled user sucecssfully : " + user);
    res.app.set('loginUser', user);
    return res.json("Enrolled user sucecssfully");
  });
});

module.exports = router;
