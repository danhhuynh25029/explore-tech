defmodule TrackerTest do
  use ExUnit.Case
  doctest Tracker

  test "greets the world" do
    assert Tracker.hello() == :world
  end
end
