(function () {

    "use strict";

    var canvas = document.getElementById("renderCanvas");
    var engine = new BABYLON.Engine(canvas, true);
    window.addEventListener("resize", function () {
        engine.resize();
    });

    var scene = new BABYLON.Scene(engine);

    scene.activeCamera = new BABYLON.FreeCamera("FreeCamera", new BABYLON.Vector3(1, 1, 1), scene);
    scene.activeCamera.attachControl(canvas);
    scene.activeCamera.keysUp.push(87);
    scene.activeCamera.keysLeft.push(65);
    scene.activeCamera.keysDown.push(83);
    scene.activeCamera.keysRight.push(68);
    scene.activeCamera.speed = 40;
    scene.activeCamera.checkCollisions = true;
    scene.activeCamera.position = new BABYLON.Vector3(500, 200, -500);
    scene.activeCamera.setTarget(new BABYLON.Vector3(-500, -200,  500));
    scene.activeCamera.attachControl(canvas, false);

    var beforeRenderFunction = function () {
        // Ensure that the camera doesn't go below ground level
        if (scene.activeCamera.position.y < 30) {
            scene.activeCamera.position.y = 30;
        }
    };
    scene.registerBeforeRender(beforeRenderFunction);

    engine.runRenderLoop(function () {
        //var tFrame = window.performance.now();
        scene.render();
        //lastTick = tFrame;
    });

    Level.init(scene);
    Client.connect(function() {}, Level.update);

})();
