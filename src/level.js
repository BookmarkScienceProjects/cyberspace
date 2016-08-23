const BABYLON = require('babylonjs');

function Update(id) {
  return {
    id,
    timestamp: 0,
    position: {x: 0, y: 0, z: 0},
    orientation: [1, 0, 0, 0],
    model: 0,
    scale: {x: 0, y: 0, z: 0},
    health: 0.0,
  };
}

function changeText(text) {
  const x = document.getElementsByClassName('content');
  let i;
  for (i = 0; i < x.length; i++) {
    x[i].innerHTML = text;
  }
}

let scene;

const objects = {};

let shadowGenerator;

const materials = {};
const models = [];

function parseInstanceInfo(inText) {
  const info = JSON.parse(inText);
  let text = `${info.Name}\n`;
  text += `${info.InstanceID}\n`;
  text += `${info.InstanceType}\n`;
  text += `CPU utilisation: ${info.CPUUtilization.toFixed(1)}%\n`;
  if (info.HasCredits) text += `CPU credits: ${info.CPUCreditBalance.toFixed(1)}\n`;
  if (info.PublicIP || info.PrivateIP) text += 'IP: ';
  if (info.PublicIP) text += `${info.PublicIP} | `;
  if (info.PrivateIP) text += info.PrivateIP;
  if (info.PublicIP || info.PrivateIP) text += '\n';
  return text;
}

function setupModels(scn) {
  // Material selection
  materials.blue = new BABYLON.StandardMaterial('texture1', scn);
  materials.blue.diffuseColor = new BABYLON.Color3(0.8, 0.8, 1);

  materials.gray = new BABYLON.StandardMaterial('texture1', scn);
  materials.gray.diffuseColor = new BABYLON.Color3(0.2, 0.9, 1);

  materials.yellow = new BABYLON.StandardMaterial('yellow', scn);
  materials.yellow.diffuseColor = new BABYLON.Color3(0.9, 0.8, 0.7);

  models[0] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
  models[0].scaling = new BABYLON.Vector3(10, 10, 10);
  models[0].isVisible = false;

  models[1] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
  models[1].scaling = new BABYLON.Vector3(30, 30, 30);
  models[1].isVisible = false;

  models[2] = BABYLON.Mesh.CreateBox('box', 1.0, scn, false, BABYLON.Mesh.DEFAULTSIDE);
  models[2].scaling = new BABYLON.Vector3(10, 10, 10);
  models[2].isVisible = false;
}

function onMeshClick(meshId) {
  const xhr = new XMLHttpRequest();
  xhr.open('POST', encodeURI('/monitor'));
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = function infoOnload() {
    if (xhr.status !== 200) {
      changeText(`Request failed. Returned status of ${xhr.status}`);
    } else {
      const text = parseInstanceInfo(xhr.responseText);
      changeText(text);
    }
  };
  xhr.send(encodeURI(`id=${meshId}`));
}

function setupOnClickAction(mesh) {
  mesh.actionManager.registerAction(
    new BABYLON.ExecuteCodeAction(
      BABYLON.ActionManager.OnPickTrigger, evt => {
        onMeshClick(evt.source.id);
      }
    )
  );
}

const updateScene = function sceneUpdater(updates) {
  Object.keys(updates).forEach((key) => {
    const update = updates[key];
    const id = update.id;

    if (typeof objects[id] === 'undefined') {
      objects[id] = models[1].clone(id);
      objects[id].id = id;
      objects[id].isVisible = true;
      objects[id].actionManager = new BABYLON.ActionManager(scene);
      shadowGenerator.getShadowMap().renderList.push(objects[id]);
      objects[id].material = new BABYLON.StandardMaterial(id, scene);
      objects[id].material.diffuseTexture = new BABYLON.Texture('/assets/square_gray.jpg', scene);
      setupOnClickAction(objects[id]);
    }

    if (update.model && objects[id].model !== update.model) {
      objects[id].model = update.model;
      //if(objects[id].material) {
        objects[id].material.dispose();
      //}
      const material = new BABYLON.StandardMaterial(id, scene);
      //material.diffuseColor = new BABYLON.Color3(0.9, 0.8, 0.7);
      //material.specularColor = new BABYLON.Color3(0.2, 0.2, 0.2);
      switch (update.model) {
        case 0:
          material.diffuseTexture = new BABYLON.Texture('/assets/square_gray.jpg', scene);
          break;
        case 1:
          material.diffuseTexture = new BABYLON.Texture('/assets/square_black.jpg', scene);
          break;
        default:
          material.diffuseTexture = new BABYLON.Texture('/assets/square_running.jpg', scene);
      }
      objects[id].material = material;
    }

    objects[id].position = update.position;
    objects[id].rotationQuaternion = new BABYLON.Quaternion(
      update.orientation[1],
      update.orientation[2],
      update.orientation[3],
      update.orientation[0]
    );
    objects[id].scaling = update.scale;

    if (objects[id].material && objects.model !== 0) {
      if (update.health > 0.99) {
        objects[id].material.emissiveColor = new BABYLON.Color3(0.0, 0.0, 0.0);
      } else if (update.health > 0.90) {
        objects[id].material.emissiveColor = new BABYLON.Color3(0.05, 0.05, 0.1);
      } else if (update.health > 0.50) {
        objects[id].material.emissiveColor = new BABYLON.Color3(0.5, 0.4, 0.3);
      } else if (update.health > 0.10) {
        objects[id].material.emissiveColor = new BABYLON.Color3(0.9, 0.6, 0.3);
      } else {
        objects[id].material.emissiveColor = new BABYLON.Color3(1, 0.0, 0.0);
      }
    }
  });
};

const entityUpdate = function entUpdate(buf) {
  const updates = {};
  let objectId;
  while (!buf.isEof()) {
    const cmd = buf.readUint8();
    switch (cmd) {
      case 1:
      {
        // INST_ENTITY_ID - we are switching the object we wish to update
        objectId = buf.readFloat32();
        updates[objectId] = new Update();
        updates[objectId].id = objectId;
        break;
      }
      case 2:
      {
        // INST_SET_POSITION
        const pos = {x: buf.readFloat32(), y: buf.readFloat32(), z: buf.readFloat32()};
        updates[objectId].position = pos;
        break;
      }
      case 3:
      {
        // INST_SET_ROTATION
        updates[objectId].orientation = [];
        updates[objectId].orientation[0] = buf.readFloat32();
        updates[objectId].orientation[1] = buf.readFloat32();
        updates[objectId].orientation[2] = buf.readFloat32();
        updates[objectId].orientation[3] = buf.readFloat32();
        break;
      }
      case 4:
      {
        // INST_SET_MODEL
        updates[objectId].model = buf.readFloat32();
        break;
      }
      case 5:
      {
        // INST_SET_SCALE
        updates[objectId].scale = {
          x: buf.readFloat32(),
          y: buf.readFloat32(),
          z: buf.readFloat32(),
        };
        break;
      }
      case 6:
      {
        updates[objectId].health = buf.readFloat32();
        break;
      }
      default:
      {
        console.log(`unknown command ${cmd}`); // eslint-disable-line
      }
    }
  }

  updateScene(updates);
};

module.exports = {

  init(s, shadow) {
    shadowGenerator = shadow;
    scene = s;
    setupModels(s);
  },

  update(buf) {
    buf.readFloat64(); // timestamp
    const msgType = buf.readUint8();
    buf.readFloat32(); // serverTick

    switch (msgType) {
      case 1:
      {
        entityUpdate(buf);
        break;
      }
      default:
      {
        console.log(`Not sure what to do with message type ${msgType}`); // eslint-disable-line
        break;
      }
    }
  },
};
