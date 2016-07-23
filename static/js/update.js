var Update = (function () {

    return function (id) {
        return {
            id: id,
            timestamp: 0,
            position: {x: 0, y: 0, z: 0},
            orientation: [1, 0, 0, 0],
            model: 0,
            scale: {x: 0, y: 0, z: 0}
        };
    }

})();
