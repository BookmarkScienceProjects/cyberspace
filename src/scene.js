const BABYLON = require('babylonjs');

const canvas = document.getElementById('renderCanvas');
const antiAlias = true;
const adaptToDeviceRation = false;

const engine = new BABYLON.Engine(canvas, antiAlias, null, adaptToDeviceRation);
window.addEventListener('resize', () => engine.resize());

BABYLON.Engine.ShadersRepository = '/assets/shaders/';

const scene = new BABYLON.Scene(engine);
//scene.clearColor = new BABYLON.Color3(0.05, 0.05, 0.05);
scene.clearColor = new BABYLON.Color3(0.2, 0.2, 0.2);
scene.ambientColor = new BABYLON.Color3(1, 1, 1);

const camera = new BABYLON.UniversalCamera('FreeCamera', new BABYLON.Vector3(1, 10, 1), scene);
//var camera = new BABYLON.WebVRFreeCamera("WVR", new BABYLON.Vector3(0, 1, -15), scene);
camera.keysUp.push(87);
camera.keysLeft.push(65);
camera.keysDown.push(83);
camera.keysRight.push(68);
camera.speed = 2;
camera.position = new BABYLON.Vector3(30, 10, 30);
camera.setTarget(new BABYLON.Vector3(0, 0, 0));
camera.attachControl(canvas, false);
scene.activeCamera = camera;

engine.runRenderLoop(() => scene.render());

module.exports = scene;
