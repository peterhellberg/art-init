const std = @import("std");

const number_of_pages = 4;

pub fn build(b: *std.Build) void {
    const exe = b.addExecutable(.{
        .name = "art-canvas",
        .root_source_file = b.path("src/canvas.zig"),
        .target = b.resolveTargetQuery(.{
            .cpu_arch = .wasm32,
            .os_tag = .freestanding,
        }),
        .optimize = .ReleaseSmall,
        .strip = true,
    });

    exe.root_module.addImport("art", b.dependency("art", .{}).module("art"));

    exe.root_module.export_symbol_names = &[_][]const u8{
        "start",
        "update",
        "draw",
        "fps",
        "offset",
        "width",
        "height",
    };

    exe.entry = .disabled;
    exe.export_memory = true;
    exe.initial_memory = std.wasm.page_size * number_of_pages;
    exe.max_memory = std.wasm.page_size * number_of_pages;
    exe.stack_size = 512;

    b.installArtifact(exe);
}
