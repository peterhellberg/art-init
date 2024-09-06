const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.resolveTargetQuery(.{
        .cpu_arch = .wasm32,
        .os_tag = .freestanding,
    });

    const exe = b.addExecutable(.{
        .name = "art-canvas",
        .root_source_file = b.path("src/canvas.zig"),
        .target = target,
        .optimize = .ReleaseSmall,
        .strip = true,
    });

    const number_of_pages = 4;

    exe.entry = .disabled;
    exe.export_memory = true;
    exe.initial_memory = std.wasm.page_size * number_of_pages;
    exe.max_memory = std.wasm.page_size * number_of_pages;
    exe.stack_size = 512;
    exe.rdynamic = true;

    b.installArtifact(exe);
}
