const BABYLON = require('babylonjs');
const scn = require('./scene.js');
const materials = require('./materials.js');

const models = [];

models[0] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[0].scaling = new BABYLON.Vector3(10, 10, 10);
models[0].isVisible = false;
models[0].material = materials.gray.clone();

models[1] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[1].scaling = new BABYLON.Vector3(10, 10, 10);
models[1].isVisible = false;
models[1].material = materials.cyan.clone();

models[2] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[2].scaling = new BABYLON.Vector3(10, 10, 10);
models[2].isVisible = false;
models[2].material = materials.sepia.clone();

module.exports = models;
