const art = @import("art");

const white = art.rgb(0xFFFFFF);
const green = art.rgb(0x7CAF3C);

var n: usize = 0;
var old: u32 = 3;
var canvas: art.Canvas(80, 45) = .{};
var color: art.RGB = white;
var pos: art.Point = .{ 0, 0 };

fn input(pad: u32) void {
    const key = art.key(pad, old);

    if (key.held(.left) and pos[0] > 0) pos[0] -|= 1;
    if (key.held(.right) and pos[0] < canvas.width - 1) pos[0] += 1;
    if (key.held(.up) and pos[1] > 0) pos[1] -|= 1;
    if (key.held(.down) and pos[1] < canvas.height - 1) pos[1] += 1;
    if (key.pressed(.z)) color = white;
    if (key.pressed(.x)) color = if (@reduce(.And, color == green)) white else green;

    old = pad;
}

fn pattern(x: i32, y: i32) art.RGBA {
    return .{
        @as(u8, @intCast(x * 3 ^ y)),
        @as(u8, @intCast(x * 3 ^ y)),
        @as(u8, @intCast(x * 3 ^ y)),
        100,
    };
}

export fn start() void {
    art.log("Hello from Zig!");
}

export fn update(pad: u32) void {
    n += 1;
    input(pad);
}

export fn draw() void {
    canvas.clear(.{ 0, 0, 0, 0 });
    canvas.fill(pattern);
    canvas.hline(14, 12, 5, .{ 255, 0, 255 });
    canvas.set(pos[0], pos[1], color);
}

export fn fps() usize {
    return 60;
}

export fn width() usize {
    return canvas.width;
}

export fn height() usize {
    return canvas.height;
}

export fn offset() [*]u8 {
    return canvas.offset();
}
