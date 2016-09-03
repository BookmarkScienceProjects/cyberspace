const BABYLON = require('babylonjs');
const light = require('./lights.js');

const shadowGenerator = new BABYLON.ShadowGenerator(1024, light);
module.exports = shadowGenerator;
