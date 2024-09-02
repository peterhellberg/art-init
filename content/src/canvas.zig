const std = @import("std");

extern fn consoleLog(arg: usize) void;

const Canvas = struct {
    const size: usize = 16;
    const fps: usize = 12;

    buffer: [size][size][4]u8 = std.mem.zeroes([size][size][4]u8),

    n: usize = 0,

    pub fn update(self: *Canvas) void {
        self.n += 1;

        self.checkerboard();

        self.set(0, 10, .{ 255, 0, 0 });
        self.set(1, 10, .{ 255, 0, 0 });
        self.set(2, 10, .{ 255, 0, 0 });
        self.set(8, 12, .{ 0, 0, @intCast(@mod(self.n, 255)) });

        const x: usize = @mod((self.n / 2) + 3, 16);

        self.set(x, 15, .{ 0, 255, 0 });

        consoleLog(x);
    }

    fn set(self: *Canvas, x: usize, y: usize, c: [3]u8) void {
        if (x < 0 or x > size or y < 0 or y > size) {
            return;
        }

        self.buffer[y][x] = .{ c[0], c[1], c[2], 255 };
    }

    fn clear(self: *Canvas, c: [3]u8) void {
        for (&self.buffer) |*row| {
            for (row) |*square| {
                square.* = .{ c[0], c[1], c[2], 255 };
            }
        }
    }

    fn checkerboard(self: *Canvas) void {
        for (&self.buffer, 0..) |*row, y| {
            for (row, 0..) |*square, x| {
                var dark = true;

                if ((y % 2) == 0) {
                    dark = false;
                }

                if ((x % 2) == 0) {
                    dark = !dark;
                }

                var p = pix(0, 0, 0);

                if (!dark) {
                    p = pix(255, 255, 255);
                }

                square.* = p;
            }
        }
    }
};

var canvas = Canvas{};

export fn canvasSize() usize {
    return Canvas.size;
}

export fn canvasFPS() usize {
    return Canvas.fps;
}

export fn canvasBufferOffset() [*]u8 {
    return @ptrCast(&canvas.buffer);
}

export fn update() void {
    canvas.update();
}

fn pix(r: u8, g: u8, b: u8) [4]u8 {
    return .{ r, g, b, 255 };
}
