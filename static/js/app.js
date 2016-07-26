(function () {

    "use strict";

    var canvas = document.getElementById("renderCanvas");
    var engine = new BABYLON.Engine(canvas, false, null, true);
    window.addEventListener("resize", function () {
        engine.resize();
    });



    var scene = new BABYLON.Scene(engine);
    //BABYLON.SceneOptimizer.OptimizeAsync(scene);
    //scene.debugLayer.show();

    scene.clearColor = new BABYLON.Color3(0.01, 0.01, 0.01);
    scene.ambientColor = new BABYLON.Color3(1, 1, 1);

    scene.activeCamera = new BABYLON.UniversalCamera("FreeCamera", new BABYLON.Vector3(1, 1, 1), scene);
    scene.activeCamera.attachControl(canvas);
    scene.activeCamera.keysUp.push(87);
    scene.activeCamera.keysLeft.push(65);
    scene.activeCamera.keysDown.push(83);
    scene.activeCamera.keysRight.push(68);
    scene.activeCamera.speed = 10;
    scene.activeCamera.checkCollisions = true;
    scene.activeCamera.position = new BABYLON.Vector3(100, 100, -500);
    scene.activeCamera.setTarget(new BABYLON.Vector3(-100, -100,  500));
    scene.activeCamera.attachControl(canvas, false);

    //new BABYLON.FxaaPostProcess("fxaa", 2.0, scene.activeCamera, 4.0);

    var mainLight = new BABYLON.HemisphericLight("Hemi0", new BABYLON.Vector3(0.3, 1, -1), scene);
    mainLight.diffuse = new BABYLON.Color3(1, 1, 1);
    mainLight.specular = new BABYLON.Color3(1, 1, 1);
    mainLight.groundColor = new BABYLON.Color3(0, 0, 0);
    mainLight.intensity = 0.95;

    // Post-process
    var blurWidth = 1.0;
    var postProcess0 = new BABYLON.PassPostProcess("Scene copy", 1.0, scene.activeCamera);
    var postProcess1 = new BABYLON.PostProcess("Down sample", "/assets/shaders/downsample", ["screenSize", "highlightThreshold"], null, 0.25, scene.activeCamera, BABYLON.Texture.BILINEAR_SAMPLINGMODE);
    postProcess1.onApply = function (effect) {
        effect.setFloat2("screenSize", postProcess1.width, postProcess1.height);
        effect.setFloat("highlightThreshold", 0.80);
    };
    var postProcess2 = new BABYLON.BlurPostProcess("Horizontal blur", new BABYLON.Vector2(1.0, 0), blurWidth, 0.25, scene.activeCamera);
    var postProcess3 = new BABYLON.BlurPostProcess("Vertical blur", new BABYLON.Vector2(0, 1.0), blurWidth, 0.25, scene.activeCamera);
    var postProcess4 = new BABYLON.PostProcess("Final compose", "/assets/shaders/compose", ["sceneIntensity", "glowIntensity", "highlightIntensity"], ["sceneSampler"], 1, scene.activeCamera);
    postProcess4.onApply = function (effect) {
        effect.setTextureFromPostProcess("sceneSampler", postProcess0);
        effect.setFloat("sceneIntensity", 1);
        effect.setFloat("glowIntensity", 0.7);
        effect.setFloat("highlightIntensity", 1.0);
    };

    //var beforeRenderFunction = function () {
    //};
    //scene.registerBeforeRender(beforeRenderFunction);

    engine.runRenderLoop(function () {
        //var tFrame = window.performance.now();
        scene.render();
        //lastTick = tFrame;
    });

    Level.init(engine, scene);
    Client.connect(function() {}, Level.update);

})();
