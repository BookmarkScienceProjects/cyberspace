const BABYLON = require('babylonjs');
const client = require('./client.js');
const level = require('./level.js');
const scene = require('./scene.js');

BABYLON.Engine.ShadersRepository = '/assets/shaders/';

// const groundMaterial = new BABYLON.StandardMaterial('ground', scene);
// groundMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
// groundMaterial.diffuseColor = new BABYLON.Color3(0.953, 0.953, 0.953);
// groundMaterial.maxSimultaneousLights = 2;

// const ground = BABYLON.Mesh.CreateGround('ground', 20000, 20000, 1, scene);
// ground.material = groundMaterial;
// ground.receiveShadows = true;

scene.registerBeforeRender(() => { scene.activeCamera.position.y = 10; });

client.connect(level.update);
