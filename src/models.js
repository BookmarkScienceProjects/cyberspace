const BABYLON = require('babylonjs');
const scn = require('./scene.js');
const materials = require('./materials.js');

const boxMat = new BABYLON.MultiMaterial('Box Multi Material', scn);
boxMat.subMaterials[0] = materials.cyan.clone();
boxMat.subMaterials[1] = materials.cyan.clone();
boxMat.subMaterials[2] = materials.darkcyan.clone();
boxMat.subMaterials[3] = materials.cyan.clone();
boxMat.subMaterials[4] = materials.cyan.clone();
boxMat.subMaterials[5] = materials.cyan.clone();

const models = [];

models[0] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[0].scaling = new BABYLON.Vector3(10, 10, 10);
models[0].isVisible = false;
models[0].material = materials.gray.clone();

models[1] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[1].scaling = new BABYLON.Vector3(10, 10, 10);
models[1].isVisible = false;
models[1].subMeshes = [];
const verticesCount = models[1].getTotalVertices();
models[1].subMeshes.push(new BABYLON.SubMesh(0, 0, verticesCount, 0, 6, models[1]));
models[1].subMeshes.push(new BABYLON.SubMesh(1, 1, verticesCount, 6, 6, models[1]));
models[1].subMeshes.push(new BABYLON.SubMesh(2, 2, verticesCount, 12, 6, models[1]));
models[1].subMeshes.push(new BABYLON.SubMesh(3, 3, verticesCount, 18, 6, models[1]));
models[1].subMeshes.push(new BABYLON.SubMesh(4, 4, verticesCount, 24, 6, models[1]));
models[1].subMeshes.push(new BABYLON.SubMesh(5, 5, verticesCount, 30, 6, models[1]));
models[1].material = boxMat.clone();

models[2] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[2].scaling = new BABYLON.Vector3(10, 10, 10);
models[2].isVisible = false;
models[2].material = materials.food.clone();

models[3] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[3].scaling = new BABYLON.Vector3(10, 10, 10);
models[3].isVisible = false;
models[3].material = materials.gray.clone();

models[4] = BABYLON.Mesh.CreateBox('grass', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
models[4].scaling = new BABYLON.Vector3(10, 10, 10);
models[4].isVisible = false;
models[4].material = materials.green.clone();

module.exports = models;
