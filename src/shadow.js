const BABYLON = require('babylonjs');
const light = require('./lights.js');

const shadowGenerator = new BABYLON.ShadowGenerator(1024*2, light);
shadowGenerator.useVarianceShadowMap = true;
module.exports = shadowGenerator;
