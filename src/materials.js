const BABYLON = require('babylonjs');
const scn = require('./scene.js');

materials.blue = new BABYLON.StandardMaterial('texture1', scn);
materials.blue.diffuseColor = new BABYLON.Color3(0.8, 0.8, 1);

materials.gray = new BABYLON.StandardMaterial('texture1', scn);
materials.gray.diffuseColor = new BABYLON.Color3(0.2, 0.9, 1);

materials.yellow = new BABYLON.StandardMaterial('yellow', scn);
materials.yellow.diffuseColor = new BABYLON.Color3(0.9, 0.8, 0.7);
