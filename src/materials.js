const BABYLON = require('babylonjs');
const scene = require('./scene.js');


const materials = {};

materials.gray = new BABYLON.StandardMaterial('gray', scene);
materials.gray.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
materials.gray.diffuseColor = new BABYLON.Color3(0.7, 0.7, 0.7);

materials.cyan = new BABYLON.StandardMaterial('purle', scene);
materials.cyan.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
materials.cyan.diffuseColor = new BABYLON.Color3(0.2, 0.9, 1);

materials.blue = new BABYLON.StandardMaterial('texture1', scene);
materials.blue.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
materials.blue.diffuseColor = new BABYLON.Color3(0.4, 0.4, 1);

materials.sepia = new BABYLON.StandardMaterial('yellow', scene);
materials.sepia.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
materials.sepia.diffuseColor = new BABYLON.Color3(0.9, 0.7, 0.5);

materials.white = new BABYLON.StandardMaterial('yellow', scene);
materials.white.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
materials.white.diffuseColor = new BABYLON.Color3(1, 1, 1);

module.exports = materials;
