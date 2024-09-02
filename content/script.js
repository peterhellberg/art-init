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

class State {
  gamepads = [0];
}

const KEY_X = 1;
const KEY_Z = 2;
const KEY_LEFT = 16;
const KEY_RIGHT = 32;
const KEY_UP = 64;
const KEY_DOWN = 128;

WebAssembly.instantiateStreaming(source, importObject).then((wasm) => {
    const state = new State();

    const mem = new Uint8Array(memory.buffer);
    const canvas = document.getElementById("art");

    const size = wasm.instance.exports.canvasSize();

    canvas.width = size;
    canvas.height = size;

    const context = canvas.getContext("2d");
    const imageData = context.createImageData(canvas.width, canvas.height);
    const bufferOffset = wasm.instance.exports.canvasBufferOffset();
    const fps = wasm.instance.exports.canvasFPS();

    let timeNextDraw = performance.now();

    const onKeyboardEvent = (event) => {
      if (event.ctrlKey || event.altKey) {
        return;
      }

      if (event.srcElement instanceof HTMLElement && event.srcElement.tagName == "INPUT") {
        return;
      }

      const down = (event.type == "keydown");

      let mask = 0;

      switch (event.code) {
      case "KeyX": case "KeyV": case "Space": case "Period":
        mask = KEY_X;
        break;
      case "KeyZ": case "KeyC": case "Comma":
        mask = KEY_Z;
        break;
      case "ArrowUp":
        mask = KEY_UP;
        break;
      case "ArrowDown":
        mask = KEY_DOWN;
        break;
      case "ArrowLeft":
        mask = KEY_LEFT;
        break;
      case "ArrowRight":
        mask = KEY_RIGHT;
        break;
      }

      if (mask != 0) {
        event.preventDefault();

        const gamepads = state.gamepads;

        if (down) {
            gamepads[0] |= mask;
        } else {
            gamepads[0] &= ~mask;
        }
      }
    };

    window.addEventListener("keydown", onKeyboardEvent);
    window.addEventListener("keyup", onKeyboardEvent);

    const frame = (timeFrameStart) => {
      requestAnimationFrame(frame);

      wasm.instance.exports.update(state.gamepads[0]);

      let calledDraw = false;

      if (timeFrameStart - timeNextDraw >= 200) {
          timeNextDraw = timeFrameStart;
      }

      while (timeFrameStart >= timeNextDraw) {
        timeNextDraw += 1000/fps;
        wasm.instance.exports.draw();
        calledDraw = true;
      }

      if (calledDraw) {
        imageData.data.set(mem.slice(
            bufferOffset,
            bufferOffset + canvas.width * canvas.height * 4
        ));

        context.putImageData(imageData, 0, 0);
      }
    };

    requestAnimationFrame(frame);
});
