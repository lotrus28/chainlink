#!/usr/bin/env node

var Web3            = require('web3'),
    contract        = require("truffle-contract"),
    RunLogJSON      = require('../../../RunLog.json');

var provider = new Web3.providers.HttpProvider("http://localhost:18545");
var RunLog = contract(RunLogJSON);
RunLog.setProvider(provider);
var devnetAddress = "0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f";

RunLog.deployed().then(function(instance) {
  return instance.request({from: devnetAddress});
}).then(console.log, console.log);
