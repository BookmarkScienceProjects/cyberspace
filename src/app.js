const BABYLON = require('babylonjs');
const Client = require('./client.js');
const Level = require('./level.js');

const canvas = document.getElementById('renderCanvas');
const antiAlias = true;
const adaptToDeviceRation = false;
const engine = new BABYLON.Engine(canvas, antiAlias, null, adaptToDeviceRation);
window.addEventListener('resize', () => engine.resize());

BABYLON.Engine.ShadersRepository = '/assets/shaders/';

const scene = new BABYLON.Scene(engine);
scene.clearColor = new BABYLON.Color3(0.05, 0.05, 0.05);
scene.ambientColor = new BABYLON.Color3(1, 1, 1);
// scene.debugLayer.show();

const ground = BABYLON.Mesh.CreateGround('ground', 20000, 20000, 1, scene);
const groundMaterial = new BABYLON.StandardMaterial('ground', scene);
groundMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
groundMaterial.diffuseColor = new BABYLON.Color3(0.2, 0.2, 0.2);
groundMaterial.maxSimultaneousLights = 2;
ground.material = groundMaterial;
ground.receiveShadows = true;


const camera = new BABYLON.UniversalCamera('FreeCamera', new BABYLON.Vector3(1, 100, 1), scene);
camera.attachControl(canvas);
camera.keysUp.push(87);
camera.keysLeft.push(65);
camera.keysDown.push(83);
camera.keysRight.push(68);
camera.speed = 20;
camera.position = new BABYLON.Vector3(100, 0, 100);
camera.setTarget(new BABYLON.Vector3(0, 0, 0));
camera.attachControl(canvas, false);
scene.activeCamera = camera;

const lightPosition = new BABYLON.Vector3(2000, 400, 2000);
const light = new BABYLON.HemisphericLight('Hemi0', lightPosition, scene);
light.intensity = 0.8;
light.diffuse = new BABYLON.Color3(1.0, 0.9, 0.9);

const mainLight = new BABYLON.PointLight('light1', lightPosition, scene);
mainLight.intensity = 0.9;
mainLight.diffuse = new BABYLON.Color3(1.0, 0.9, 0.85);
mainLight.specular = new BABYLON.Color3(1, 1, 1);
mainLight.groundColor = new BABYLON.Color3(0.2, 0.2, 0.2);

const shadowGenerator = new BABYLON.ShadowGenerator(1024, mainLight);

// Post-process
const blurWidth = 1;
const postProcess0 = new BABYLON.PassPostProcess('Scene copy', 1.0, scene.activeCamera);
const postProcess1 = new BABYLON.PostProcess(
  'Down sample',
  'downsample',
  ['screenSize', 'highlightThreshold'],
  null,
  0.25,
  scene.activeCamera,
  BABYLON.Texture.BILINEAR_SAMPLINGMODE
);

postProcess1.onApply = function pp1OnApply(effect) {
  effect.setFloat2('screenSize', postProcess1.width, postProcess1.height);
  effect.setFloat('highlightThreshold', 0.80);
};

const postProcess2 = new BABYLON.BlurPostProcess('Horizontal blur', new BABYLON.Vector2(1.0, 0), blurWidth, 0.25, scene.activeCamera);  // eslint-disable-line
const postProcess3 = new BABYLON.BlurPostProcess('Vertical blur', new BABYLON.Vector2(0, 1.0), blurWidth, 0.25, scene.activeCamera);  // eslint-disable-line
const postProcess4 = new BABYLON.PostProcess('Final compose', '/assets/shaders/compose', ['sceneIntensity', 'glowIntensity', 'highlightIntensity'], ['sceneSampler'], 1, scene.activeCamera);  // eslint-disable-line
postProcess4.onApply = function ps4OnApply(effect) {
  effect.setTextureFromPostProcess('sceneSampler', postProcess0);
  effect.setFloat('sceneIntensity', 0.9);
  effect.setFloat('glowIntensity', 0.3);
  effect.setFloat('highlightIntensity', 1.0);
};

function beforeRenderFunction() {
  scene.activeCamera.position.y = 300;
}

scene.registerBeforeRender(beforeRenderFunction);

engine.runRenderLoop(() => scene.render());

Level.init(scene, shadowGenerator);

Client.connect(Level.update);
