let Oracle = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/Oracle.sol");
let RunLog = artifacts.require("./RunLog.sol");

module.exports = function(deployer) {
  deployer.deploy(RunLog, Oracle.address);
};
