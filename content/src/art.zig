const std = @import("std");

pub const size: usize = 32;
pub const fps: usize = 60;

pub const Key = struct {
    x: bool,
    z: bool,
    left: bool,
    right: bool,
    up: bool,
    down: bool,
};

pub const KeyX = 1;
pub const KeyZ = 2;
pub const KeyLeft = 16;
pub const KeyRight = 32;
pub const KeyUp = 64;
pub const KeyDown = 128;

pub fn key(pad: u32) Key {
    return .{
        .x = (pad & KeyX != 0),
        .z = (pad & KeyZ != 0),
        .left = (pad & KeyLeft != 0),
        .right = (pad & KeyRight != 0),
        .up = (pad & KeyUp != 0),
        .down = (pad & KeyDown != 0),
    };
}

pub const Size = @Vector(2, usize);
pub const RGB = @Vector(3, u8);
pub const RGBA = @Vector(4, u8);

pub var buffer: [size][size][4]u8 = std.mem.zeroes([size][size][4]u8);

pub fn clear(c: RGBA) void {
    for (&buffer) |*row| {
        for (row) |*square| {
            square.* = c;
        }
    }
}

pub fn fill(color: *const fn (x: usize, y: usize) RGBA) void {
    for (&buffer, 0..) |*row, y| {
        for (row, 0..) |*square, x| {
            const c = color(x, y);

            if (c[3] > 0) square.* = c;
        }
    }
}

pub fn set(x: usize, y: usize, c: RGB) void {
    if (x > size or y > size) {
        return;
    }

    buffer[y][x] = .{ c[0], c[1], c[2], 255 };
}

pub fn hline(x: usize, y: usize, w: usize, color: RGB) void {
    var to = x + w;

    if (to >= size - 1) {
        to = size;
    }

    for (x..to) |rx| {
        set(rx, y, color);
    }
}

pub fn vline(x: usize, y: usize, h: usize, color: RGB) void {
    var to = y + h;

    if (to >= size - 1) {
        to = size;
    }

    for (y..to) |ry| {
        set(x, ry, color);
    }
}

pub const rectArgs = struct {
    size: Size = .{ 1, 1 },
    color: RGB = .{ 255, 255, 255 },
};

pub fn rect(x: usize, y: usize, args: rectArgs) void {
    for (x..x + args.size[0]) |rx| {
        for (y..y + args.size[1]) |ry| {
            set(rx, ry, args.color);
        }
    }
}

pub const boxArgs = struct {
    size: Size = .{ 1, 1 },
    color: RGB = .{ 0, 0, 0 },
    fill: bool = false,
    fillColor: RGB = .{ 0, 0, 0 },
};

pub fn box(x: usize, y: usize, args: boxArgs) void {
    if (args.fill) {
        rect(x, y, .{
            .size = args.size -| Size{ 1, 1 },
            .color = args.fillColor,
        });
    }

    const color = args.color;
    const w = args.size[0];
    const h = args.size[1];

    const by = y + (h -| 1);
    const rx = x + (w -| 1);

    if (x > size or y > size) return;

    hline(x, y, w, color);

    if (by > size) return;

    hline(x, by, w, color);
    vline(x, y + 1, h -| 2, color);

    if (rx < size) vline(rx, y + 1, h -| 2, color);
}

pub fn rgb(hex: u32) RGB {
    return .{
        @intCast(hex >> 16 & 0xFF),
        @intCast(hex >> 8 & 0xFF),
        @intCast(hex & 0xFF),
    };
}

pub fn rgba(hexa: u32) RGBA {
    return .{
        @intCast(hexa >> 32 & 0xFF),
        @intCast(hexa >> 16 & 0xFF),
        @intCast(hexa >> 8 & 0xFF),
        @intCast(hexa & 0xFF),
    };
}

pub extern fn consoleLog(arg: usize) void;

export fn canvasSize() usize {
    return size;
}

export fn canvasFPS() usize {
    return fps;
}

export fn canvasBufferOffset() [*]u8 {
    return @ptrCast(&buffer);
}
