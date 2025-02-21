// GL
extern fn glClearColor(r: f32, g: f32, b: f32, a: f32) void;
extern fn glEnable(cap: c_uint) void;
extern fn glDepthFunc(func: c_uint) void;
extern fn glClear(mask: c_uint) void;
extern fn glGetAttribLocation(program: c_uint, start: *const u8, len: c_uint) c_int;
extern fn glGetUniformLocation(program: c_uint, start: *const u8, len: c_uint) c_int;
extern fn glUniform1f(location: c_int, v0: f32) void;
extern fn glUniform2fv(location: c_int, v0: f32, v1: f32) void;
extern fn glUniform4fv(location: c_int, x: f32, y: f32, z: f32, w: f32) void;
extern fn glCreateBuffer() c_uint;
extern fn glBindBuffer(target: c_uint, buffer: c_uint) void;
extern fn glBufferData(target: c_uint, size: *const f32, data: c_uint, usage: c_uint) void;
extern fn glUseProgram(program: c_uint) void;
extern fn glEnableVertexAttribArray(index: c_uint) void;
extern fn glVertexAttribPointer(index: c_uint, size: c_uint, typ: c_uint, normalized: c_uint, stride: c_uint, pointer: c_uint) void;
extern fn glDrawArrays(mode: c_uint, first: c_uint, count: c_uint) void;

// Identifier constants pulled from WebGLRenderingContext
const GL_VERTEX_SHADER: c_uint = 35633;
const GL_FRAGMENT_SHADER: c_uint = 35632;
const GL_ARRAY_BUFFER: c_uint = 34962;
const GL_TRIANGLES: c_uint = 4;
const GL_STATIC_DRAW: c_uint = 35044;
const GL_f32: c_uint = 5126;
const GL_DEPTH_TEST: c_uint = 2929;
const GL_LEQUAL: c_uint = 515;
const GL_COLOR_BUFFER_BIT: c_uint = 16384;
const GL_DEPTH_BUFFER_BIT: c_uint = 256;

// Math
extern fn sin(f64) f64; // returns the sine of a number in radians.
extern fn cos(f64) f64; // returns the cosine of a number in radians.

// Shaders
extern fn compileShader(source: *const u8, len: c_uint, type: c_uint) c_uint;
extern fn linkShaderProgram(vertexShaderId: c_uint, fragmentShaderId: c_uint) c_uint;

var program: c_uint = undefined;
var mouseUniformLocation: c_int = undefined;
var offsetUniformLocation: c_int = undefined;
var positionAttributeLocation: c_int = undefined;
var resolutionUniformLocation: c_int = undefined;
var timeUniformLocation: c_int = undefined;
var positionBuffer: c_uint = undefined;

export fn onInit() void {
    // GL setup
    glEnable(GL_DEPTH_TEST);
    glDepthFunc(GL_LEQUAL);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    glClearColor(0.1, 0.1, 0.1, 1.0);

    program = compileAndLinkShaders();

    // Uniforms
    const u_mouse = "u_mouse";
    const u_offset = "u_offset";
    const u_resolution = "u_resolution";
    const u_time = "u_time";

    mouseUniformLocation = glGetUniformLocation(program, &u_mouse[0], u_mouse.len);
    offsetUniformLocation = glGetUniformLocation(program, &u_offset[0], u_offset.len);
    resolutionUniformLocation = glGetUniformLocation(program, &u_resolution[0], u_resolution.len);
    timeUniformLocation = glGetUniformLocation(program, &u_time[0], u_time.len);

    // Buffer
    positionBuffer = glCreateBuffer();

    glBindBuffer(GL_ARRAY_BUFFER, positionBuffer);
    glBufferData(GL_ARRAY_BUFFER, &triangles[0], triangles.len, GL_STATIC_DRAW);
}

const triangles = [_]f32{
    -1, 1, 1, -1, 1,  1,
    -1, 1, 1, -1, -1, -1,
};

var p: c_int = 0;
var t: f32 = 0;
var x: f32 = 0;

export fn onAnimationFrame(ts: c_int, w: f32, h: f32, mx: f32, my: f32) void {
    t = @as(f32, @floatFromInt(ts)) / 1000;

    x += nx(ts) / 2000;

    if (x > 2) x = -2; // wrap around

    glUseProgram(program);

    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 2, GL_f32, 0, 0, 0);
    glBindBuffer(GL_ARRAY_BUFFER, positionBuffer);

    const y: f32 = @floatCast(sin(cos(x * 5) * 0.3));

    glUniform2fv(mouseUniformLocation, mx, my);
    glUniform2fv(resolutionUniformLocation, w, h);
    glUniform1f(timeUniformLocation, t);

    glUniform4fv(offsetUniformLocation, x, y, 0, 0);
    glDrawArrays(GL_TRIANGLES, 0, triangles.len / 2);
    glUniform4fv(offsetUniformLocation, -x, y, 0, 0);
    glDrawArrays(GL_TRIANGLES, 0, triangles.len / 2);

    p = ts;
}

fn nx(ts: c_int) f32 {
    const delta = if (p > 0) ts - p else 0;

    return @as(f32, @floatFromInt(delta));
}

fn compileAndLinkShaders() c_uint {
    const vertex = @embedFile("vertex.glslx");
    const fragment = @embedFile("fragment.glslx");

    return linkShaderProgram(
        compileShader(&vertex[0], vertex.len, GL_VERTEX_SHADER),
        compileShader(&fragment[0], fragment.len, GL_FRAGMENT_SHADER),
    );
}
