const std = @import("std");
const expert = std.testing.expect;

test "success" {
    try expert(true);
}
