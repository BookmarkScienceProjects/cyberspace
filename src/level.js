const BABYLON = require('babylonjs');

const scene = require('./scene.js');
//const materials = require('./materials.js');
const models = require('./models.js');
const shadowGenerator = require('./shadow.js');

function Update(id) {
  return {
    id,
    timestamp: 0,
    position: { x: 0, y: 0, z: 0 },
    orientation: [1, 0, 0, 0],
    model: 0,
    scale: { x: 0, y: 0, z: 0 },
    health: 0.0
  };
}

function changeText(text) {
  const x = document.getElementsByClassName('content');
  let i;
  for (i = 0; i < x.length; i++) {
    x[i].innerHTML = text;
  }
}

const objects = {};

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

function onMeshClick(meshId) {
  const xhr = new XMLHttpRequest();
  xhr.open('POST', encodeURI('/click'));
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

    if (update.state === 0) {
      return;
    }

    if (typeof objects[id] === 'undefined') {
      if (typeof update.model === 'number') {
        objects[id] = models[update.model].clone(id);
      } else {
        objects[id] = models[0].clone(id);
      }
      objects[id].id = id;
      objects[id].isVisible = true;
      objects[id].actionManager = new BABYLON.ActionManager(scene);
      shadowGenerator.getShadowMap().renderList.push(objects[id]);
      setupOnClickAction(objects[id]);
    }

    objects[id].position = update.position;
    objects[id].rotationQuaternion = new BABYLON.Quaternion(
      update.orientation[1],
      update.orientation[2],
      update.orientation[3],
      update.orientation[0]
    );
    objects[id].scaling = update.scale;
    // we are going to ignore any height offset just because I don't want to deal with the
    // steering output in the backen
    objects[id].position.y = objects[id].scaling.y/2;
  });
};

const entityRemove = function entityRemove(buf) {
  while (!buf.isEof()) {
    const cmd = buf.readUint8();
    switch (cmd) {
      case 1: {
        // INST_ENTITY_ID - we are switching the object we wish to update
        const objectId = buf.readFloat32();
        // check if object exists before disposing
        if(objects[objectId] !== undefined) {
          objects[objectId].dispose();
        }
        break;
      }

      default: {
        console.log(`unknown command ${cmd}`); // eslint-disable-line
      }
    }
  }

};

const entityUpdate = function entUpdate(buf) {
  const updates = {};
  let objectId;
  while (!buf.isEof()) {
    const cmd = buf.readUint8();
    switch (cmd) {
      case 1: {
        // INST_ENTITY_ID - we are switching the object we wish to update
        objectId = buf.readFloat32();
        updates[objectId] = new Update();
        updates[objectId].id = objectId;
        break;
      }
      case 2: {
        // INST_SET_POSITION
        updates[objectId].position = {
          x: buf.readFloat32(),
          y: buf.readFloat32(),
          z: buf.readFloat32()
        };
        break;
      }
      case 3: {
        // INST_SET_ROTATION
        updates[objectId].orientation = [];
        updates[objectId].orientation[0] = buf.readFloat32();
        updates[objectId].orientation[1] = buf.readFloat32();
        updates[objectId].orientation[2] = buf.readFloat32();
        updates[objectId].orientation[3] = buf.readFloat32();
        break;
      }
      case 4: {
        // INST_SET_MODEL
        updates[objectId].model = buf.readInt32();
        break;
      }
      case 5: {
        // INST_SET_SCALE
        updates[objectId].scale = {
          x: buf.readFloat32(),
          y: buf.readFloat32(),
          z: buf.readFloat32(),
        };
        break;
      }
      case 6: {
        updates[objectId].state = buf.readFloat32();
        break;
      }
      default: {
        console.log(`unknown command ${cmd}`); // eslint-disable-line
      }
    }
  }

  updateScene(updates);
};

module.exports = {

  update(buf) {
    buf.readFloat64(); // timestamp
    const msgType = buf.readUint8();
    buf.readFloat32(); // serverTick
    switch (msgType) {
      case 1: {
        entityUpdate(buf);
        break;
      }
      case 2: {
        entityRemove(buf);
        break;
      }
      default: {
        console.log(`Not sure what to do with message type ${msgType}`); // eslint-disable-line
        break;
      }
    }
  },
};
