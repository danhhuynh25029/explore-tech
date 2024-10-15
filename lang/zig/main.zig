const std = @import("std");

pub fn main() void {
    std.debug.print("Hello, {s}!\n", .{"World"});

    const a = [5]u8{ 'h', 'e', 'l', 'l', 'o' };
    const b = [_]u8{ 'w', 'o', 'r', 'l', 'd', '1' };

    std.debug.print("Hello, {d}!\n", .{a.len});

    std.debug.print("Hello, {d}!\n", .{b.len});

    // var i: u8 = 1;

    // while (i < 100) {
    //     std.debug.print("i : {d}\n", .{i});
    //     i *= 2;
    // }

    for (a, 0..) |character, index| {
        std.debug.print("character : {c} , index : {d}\n", .{ character, index });
    }
}
