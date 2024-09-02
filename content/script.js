/* number of pages, specified in build.zig */
var memory = new WebAssembly.Memory({
  initial: 4,
  maximum: 4,
});

var importObject = {
    env: {
        consoleLog: (arg) => console.log(arg), // Useful for debugging on Zig's side
        memory: memory,
    },
};

var source = fetch("zig-out/bin/art-canvas.wasm");

WebAssembly.instantiateStreaming(source, importObject).then((wasm) => {
    const mem = new Uint8Array(memory.buffer);
    const size = wasm.instance.exports.canvasSize();
    const canvas = document.getElementById("art");

    // Set the canvas width and height to a square of the size retrieved from WASM
    canvas.width = size;
    canvas.height = size;

    // Get a 2D context to draw into
    const context = canvas.getContext("2d");

    // Create image data for the entire canvas
    const imageData = context.createImageData(canvas.width, canvas.height);

    // Get the buffer offset for the canvas
    const bufferOffset = wasm.instance.exports.canvasBufferOffset();

    // Get the desired FPS
    const fps = wasm.instance.exports.canvasFPS();

    // When we should perform the next update
    let timeNextUpdate = performance.now();

    const draw = (timeFrameStart) => {
      requestAnimationFrame(draw);

      let calledUpdate = false;

      // Prevent timeFrameStart from getting too far ahead and death spiralling
      if (timeFrameStart - timeNextUpdate >= 200) {
          timeNextUpdate = timeFrameStart;
      }

      while (timeFrameStart >= timeNextUpdate) {
        timeNextUpdate += 1000/fps;

        wasm.instance.exports.update();

        calledUpdate = true;
      }

      if (calledUpdate) {
        // Set the image data from a slice of memory from bufferOffset
        imageData.data.set(mem.slice(
            bufferOffset,
            bufferOffset + canvas.width * canvas.height * 4
        ));

        // Draw the image data into the canvas
        context.putImageData(imageData, 0, 0);
      }
    };

    requestAnimationFrame(draw);
});
