(function () {

    "use strict";

    var canvas = document.getElementById("renderCanvas");
    var antialias = true;
    var adaptToDeviceRation = false;
    var engine = new BABYLON.Engine(canvas, antialias, null, adaptToDeviceRation);
    window.addEventListener("resize", function () {
        engine.resize();
    });

    BABYLON.Engine.ShadersRepository = "/assets/shaders/";

    var scene = new BABYLON.Scene(engine);
    scene.clearColor = new BABYLON.Color3(0.15, 0.13, 0.13);
    scene.ambientColor = new BABYLON.Color3(1, 1, 1);
    //scene.debugLayer.show();

    var ground = BABYLON.Mesh.CreateGround("ground", 20000, 20000, 1, scene);
    var groundMaterial = new BABYLON.StandardMaterial("ground", scene);
    groundMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
    groundMaterial.diffuseColor = new BABYLON.Color3(0.2, 0.2, 0.2);
    groundMaterial.maxSimultaneousLights = 2;
    ground.material = groundMaterial;
    ground.receiveShadows = true;

    var camera = new BABYLON.UniversalCamera("FreeCamera", new BABYLON.Vector3(1, 100, 1), scene);
    camera.attachControl(canvas);
    camera.keysUp.push(87);
    camera.keysLeft.push(65);
    camera.keysDown.push(83);
    camera.keysRight.push(68);
    camera.speed = 10;
    camera.position = new BABYLON.Vector3(500, 300, -500);
    camera.setTarget(new BABYLON.Vector3(0, 0,  0));
    camera.attachControl(canvas, false);
    scene.activeCamera = camera;

    var lightPosition = new BABYLON.Vector3(600, 300, 600);
    var light = new BABYLON.HemisphericLight("Hemi0", lightPosition, scene);
    light.intensity = 0.9;
    light.diffuse = new BABYLON.Color3(1.0, 0.9, 0.9);

    var mainLight  = new BABYLON.PointLight("light1", lightPosition, scene);
    mainLight.intensity = 0.7;
    mainLight.diffuse = new BABYLON.Color3(1.0, 0.9, 0.9);
    mainLight.specular = new BABYLON.Color3(1, 1, 1);
    mainLight.groundColor = new BABYLON.Color3(0.2, 0.2, 0.2);
    mainLight.intensity = 0.4;
    var sphere = BABYLON.Mesh.CreateSphere("sphere", 5.0, 5.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
    sphere.position = mainLight.position;
    sphere.material = new BABYLON.StandardMaterial("sphere", scene);
    sphere.material.emissiveColor = mainLight.diffuse;

    // Post-process
    var blurWidth = 1;
    var postProcess0 = new BABYLON.PassPostProcess("Scene copy", 1.0, scene.activeCamera);
    var postProcess1 = new BABYLON.PostProcess("Down sample", "downsample", ["screenSize", "highlightThreshold"], null, 0.25, scene.activeCamera, BABYLON.Texture.BILINEAR_SAMPLINGMODE);
    postProcess1.onApply = function (effect) {
        effect.setFloat2("screenSize", postProcess1.width, postProcess1.height);
        effect.setFloat("highlightThreshold", 0.80);
    };
    var postProcess2 = new BABYLON.BlurPostProcess("Horizontal blur", new BABYLON.Vector2(1.0, 0), blurWidth, 0.25, scene.activeCamera);
    var postProcess3 = new BABYLON.BlurPostProcess("Vertical blur", new BABYLON.Vector2(0, 1.0), blurWidth, 0.25, scene.activeCamera);
    var postProcess4 = new BABYLON.PostProcess("Final compose", "/assets/shaders/compose", ["sceneIntensity", "glowIntensity", "highlightIntensity"], ["sceneSampler"], 1, scene.activeCamera);
    postProcess4.onApply = function (effect) {
        effect.setTextureFromPostProcess("sceneSampler", postProcess0);
        effect.setFloat("sceneIntensity", 0.9);
        effect.setFloat("glowIntensity", 0.3);
        effect.setFloat("highlightIntensity", 1.0);
    };

    var beforeRenderFunction = function () {
        scene.activeCamera.position.y = 300;
    };
    scene.registerBeforeRender(beforeRenderFunction);

    engine.runRenderLoop(function () {
        //var tFrame = window.performance.now();
        scene.render();
        //lastTick = tFrame;
    });

    Level.init(engine, scene, mainLight);
    Client.connect(function() {}, Level.update);

})();
