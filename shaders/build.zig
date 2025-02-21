const std = @import("std");

pub fn build(b: *std.Build) void {
    const exe = b.addExecutable(.{
        .name = "art-canvas",
        .root_source_file = b.path("src/webgl.zig"),
        .target = b.resolveTargetQuery(.{
            .cpu_arch = .wasm32,
            .os_tag = .freestanding,
        }),
        .optimize = .ReleaseSmall,
        .strip = true,
    });

    exe.entry = .disabled;
    exe.rdynamic = true;

    b.installArtifact(exe);
}
