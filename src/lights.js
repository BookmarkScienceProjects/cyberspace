const BABYLON = require('babylonjs');

const scene = require('./scene.js');

const lightPosition = new BABYLON.Vector3(200, 400, 200);
const light = new BABYLON.HemisphericLight('Hemi0', new BABYLON.Vector3(2, 4, 2).normalize(), scene);
light.intensity = 0.7;
light.diffuse = new BABYLON.Color3(1.0, 0.9, 0.9);
light.groundColor = new BABYLON.Color3(0.5, 0.5, 0.5);

const mainLight = new BABYLON.DirectionalLight("dir", new BABYLON.Vector3(-2, -4, -2).normalize(), scene);
mainLight.position = lightPosition;
mainLight.intensity = 0.5;
mainLight.diffuse = new BABYLON.Color3(1.0, 0.9, 0.85);
mainLight.specular = new BABYLON.Color3(1, 1, 1);
mainLight.groundColor = new BABYLON.Color3(0.5, 0.5, 0.5);

module.exports = mainLight;
