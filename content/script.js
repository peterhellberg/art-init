const binary = "zig-out/bin/art-canvas.wasm";

const state = {};

fetch(binary).then((source) => {
  WebAssembly.instantiateStreaming(source, {
    env: {
      Log: (ptr, len) => {
        // Useful for debugging on Zig's side
        const buf = state.mem.slice(ptr, ptr+len);

        console.log(new TextDecoder().decode(buf));
      }
    },
  }).then((wasm) => {
    state.mem = new Uint8Array(wasm.instance.exports.memory.buffer);

    const fps = wasm.instance.exports.fps();
    const art = document.getElementById("art");

    art.width = wasm.instance.exports.width();
    art.height = wasm.instance.exports.height();

    const ctx = art.getContext("2d");
    const image = ctx.createImageData(art.width, art.height);

    const offset = wasm.instance.exports.offset();
    const offsetEnd = offset + art.width * art.height * 4;

    wasm.instance.exports.start();

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

      const KEY_X = 1;
      const KEY_Z = 2;
      const KEY_LEFT = 16;
      const KEY_RIGHT = 32;
      const KEY_UP = 64;
      const KEY_DOWN = 128;

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

        if (down) {
          state.gamepad |= mask;
        } else {
          state.gamepad &= ~mask;
        }
      }
    };

    window.addEventListener("keydown", onKeyboardEvent);
    window.addEventListener("keyup", onKeyboardEvent);

    const frame = (timeFrameStart) => {
      requestAnimationFrame(frame);

      let calledDraw = false;

      if (timeFrameStart - timeNextDraw >= 200) {
        timeNextDraw = timeFrameStart;
      }

      while (timeFrameStart >= timeNextDraw) {
        timeNextDraw += 1000/fps;
        wasm.instance.exports.update(state.gamepad);
        wasm.instance.exports.draw();
        calledDraw = true;
      }

      if (calledDraw) {
        image.data.set(state.mem.slice(offset, offsetEnd));
        ctx.putImageData(image, 0, 0);
      }
    };

    if(fps == 0) {
      wasm.instance.exports.draw();
      image.data.set(state.mem.slice(offset, offsetEnd));
      ctx.putImageData(image, 0, 0);
    } else {
      requestAnimationFrame(frame);
    }
  });
});
