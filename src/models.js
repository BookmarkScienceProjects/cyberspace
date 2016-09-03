const BABYLON = require('babylonjs');
const scn = require('./scene.js');

const models = [];

models[0] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[0].scaling = new BABYLON.Vector3(10, 10, 10);
models[0].isVisible = false;

models[1] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[1].scaling = new BABYLON.Vector3(30, 30, 30);
models[1].isVisible = false;

models[2] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[2].scaling = new BABYLON.Vector3(10, 10, 10);
models[2].isVisible = false;


module.exports = models;
