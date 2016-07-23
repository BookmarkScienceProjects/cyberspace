var Level = (function () {

    "use strict";

    var scene;
    // contains the last unix timestamp received from the websocket
    var timestamp;

    var serverTick;
    var objects = {};

    var materials = {};
    var models = [];

    function updateScene(updates) {
        for (var id in updates) {
            if (!updates.hasOwnProperty(id)) {
                continue;
            }

            // entity needs to be created
            if (typeof objects[id] === 'undefined') {
                objects[id] = models[updates[id].model].createInstance(id);
                objects[id].isVisible = true;
            }

            objects[id].position = updates[id].position;
            objects[id].rotationQuaternion = new BABYLON.Quaternion(updates[id].orientation[1], updates[id].orientation[2], updates[id].orientation[3], updates[id].orientation[0]);
            objects[id].scaling = updates[id].scale;

        }
    }

    return {

        init: function (s) {
            scene = s;

            scene.clearColor = new BABYLON.Color3(0.1, 0.1, 0.1);
            scene.ambientColor = new BABYLON.Color3(1, 1, 1);

            // Material selection
            materials.blue = new BABYLON.StandardMaterial("texture1", scene);
            materials.blue.diffuseColor = new BABYLON.Color3(0.0, 0.0, 0.4);

            materials.gray = new BABYLON.StandardMaterial("texture1", scene);
            materials.gray.diffuseColor = new BABYLON.Color3(0.2, 0.2, 0.2);

            materials.red = new BABYLON.StandardMaterial("red", scene);
            materials.red.diffuseColor = new BABYLON.Color3(0.9, 0.2, 0.2);

            materials.pink = new BABYLON.StandardMaterial("texture1", scene);
            materials.pink.diffuseColor = new BABYLON.Color3(1.0, 0.2, 0.7);

            materials.moccasin = new BABYLON.StandardMaterial("texture1", scene);
            materials.moccasin.diffuseColor = new BABYLON.Color3(1.0,0.9, 0.8);

            materials.lightPink = new BABYLON.StandardMaterial("texture1", scene);
            materials.lightPink.diffuseColor = new BABYLON.Color3(.7, 0.3, .7);

            materials.yellow = new BABYLON.StandardMaterial("yellow", scene);
            materials.yellow.diffuseColor = new BABYLON.Color3(0.9, 0.8, 0.7);

            materials.green = new BABYLON.StandardMaterial("texture1", scene);
            materials.green.diffuseColor = new BABYLON.Color3(.5, 1.0, .4);

            models[0] = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
            models[0].scaling = new BABYLON.Vector3(10, 10, 10);
            models[0].material = materials.gray;
            models[0].isVisible = false;

            models[1] = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
            models[1].scaling = new BABYLON.Vector3(30, 30, 30);
            models[1].isVisible = false;
            var multi=new BABYLON.MultiMaterial("nuggetman",scene);
            //multi.subMaterials.push(green);
            multi.subMaterials.push(materials.yellow);
            multi.subMaterials.push(materials.yellow);
            //multi.subMaterials.push(materials.red);
            multi.subMaterials.push(materials.yellow);
            multi.subMaterials.push(materials.yellow);
            multi.subMaterials.push(materials.yellow);
            multi.subMaterials.push(materials.yellow);

            models[1].subMeshes=[];
            var verticesCount=models[1].getTotalVertices();
            models[1].subMeshes.push(new BABYLON.SubMesh(0, 0, verticesCount, 0, 6, models[1]));
            models[1].subMeshes.push(new BABYLON.SubMesh(1, 1, verticesCount, 6, 6, models[1]));
            models[1].subMeshes.push(new BABYLON.SubMesh(2, 2, verticesCount, 12, 6, models[1]));
            models[1].subMeshes.push(new BABYLON.SubMesh(3, 3, verticesCount, 18, 6, models[1]));
            models[1].subMeshes.push(new BABYLON.SubMesh(4, 4, verticesCount, 24, 6, models[1]));
            models[1].subMeshes.push(new BABYLON.SubMesh(5, 5, verticesCount, 30, 6, models[1]));
            models[1].material=multi;

            var mainLight = new BABYLON.HemisphericLight("Hemi0", new BABYLON.Vector3(0.3, 1, -1), scene);
            mainLight.diffuse = new BABYLON.Color3(1,1,1);
            mainLight.specular = new BABYLON.Color3(1, 1, 1);
            mainLight.groundColor = new BABYLON.Color3(0, 0, 0);
            mainLight.intensity = 0.95;


            var canvas = new BABYLON.ScreenSpaceCanvas2D(scene, {
                id: "ScreenCanvas",
                size: new BABYLON.Size(300, 100),
                backgroundFill: "#4040408F",
                backgroundRoundRadius: 0,
                children: [
                    new BABYLON.Text2D("cyberspace", {
                        id: "text",
                        marginAlignment: "h: center, v:center",
                        fontName: "20pt Arial",
                    })
                ]
            });

        },

        update: function (buf) {

            timestamp = buf.readFloat64();
            var msgType = buf.readUint8();
            serverTick = buf.readFloat32();

            var objectId;
            var updates = [];

            while (!buf.isEof()) {
                var cmd = buf.readUint8();
                switch (cmd) {
                    // INST_ENTITY_ID - we are switching the object we wish to update
                    case 1:
                        objectId = buf.readFloat32();
                        updates[objectId] = new Update();
                        break;
                    // INST_SET_POSITION
                    case 2:
                        var pos = {x: buf.readFloat32(), y: buf.readFloat32(), z: buf.readFloat32()};
                        updates[objectId].position = pos;
                        break;
                    // INST_SET_ROTATION
                    case 3:
                        updates[objectId].orientation = [];
                        updates[objectId].orientation[0] = buf.readFloat32();
                        updates[objectId].orientation[1] = buf.readFloat32();
                        updates[objectId].orientation[2] = buf.readFloat32();
                        updates[objectId].orientation[3] = buf.readFloat32();
                        break;
                    // INST_SET_MODEL
                    case 4:
                        updates[objectId].model = buf.readFloat32();
                        break;
                    // INST_SET_SCALE
                    case 5:
                        updates[objectId].scale = {x: buf.readFloat32(), y: buf.readFloat32(), z: buf.readFloat32()};
                        break;
                    // INST_KILL - remove this entrity
                    case 6:
                        break;
                }
            }

            updateScene(updates);
        }
    }
})();
