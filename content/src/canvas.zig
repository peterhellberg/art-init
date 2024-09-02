const art = @import("art.zig");

fn pattern(x: usize, y: usize) art.RGBA {
    return .{
        @as(u8, @intCast(x * 2 ^ y)),
        @as(u8, @intCast(x ^ y)),
        @as(u8, @intCast(x ^ y)),
        @intCast(@mod(state.n * 1, 255)),
    };
}

const orange = art.rgb(0xFF6600);
const white = art.rgb(0xFFFFFF);

var state: struct {
    n: usize = 0,
    old: u32 = 0,
    pos: art.Size = .{ 0, 0 },
    color: art.RGB = white,

    const State = @This();

    fn draw(self: *State) void {
        art.clear(.{ 0, 0, 0, 255 });

        art.fill(pattern);

        art.box(11, 2, .{
            .size = .{ 6, 9 },
            .color = .{ 255, 0, 0 },
            .fill = true,
            .fillColor = .{ 0, 255, 0 },
        });

        art.set(self.pos[0], self.pos[1], self.color);

        art.vline(14, 12, 5, .{ 255, 0, 255 });
    }

    fn update(self: *State, pad: u32) void {
        self.n += 1;

        const key = art.key(pad);
        const old = art.key(self.old);

        if (key.x and !old.x) {
            self.color = if (@reduce(.And, self.color == orange)) white else orange;
        }

        if (key.z) self.color = white;

        if (key.up and !old.up) {
            self.pos[1] -|= 1;
        }

        if (key.down and !old.down) if (self.pos[1] < art.size - 1) {
            self.pos[1] += 1;
        };

        if (key.left and !old.left) {
            self.pos[0] -|= 1;
        }

        if (key.right and !old.right) {
            if (self.pos[0] < art.size - 1) {
                self.pos[0] += 1;
            }
        }

        self.old = pad;
    }
} = .{};

export fn draw() void {
    state.draw();
}

export fn update(pad: u32) void {
    state.update(pad);
}
