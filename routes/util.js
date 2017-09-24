var fs = require('fs');

exports.setCcConfig = function(chaincodeName) {
  try {
    return JSON.parse(fs.readFileSync(__dirname + '/../config/'+ chaincodeName + '.json', 'utf8'));
  } catch (err) {
    console.log(chaincodeName + ".json is missing or invalid file, Rerun the program with right file");
    process.exit();
  }
}

exports.setCredentials = function() {
  try {
    network = JSON.parse(fs.readFileSync(__dirname + '/../ServiceCredentials.json', 'utf8'));
    if (network.credentials) return network.credentials;
  } catch (err) {
    console.log("ServiceCredentials.json is missing or invalid file, Rerun the program with right file")
    process.exit();
  }
}
