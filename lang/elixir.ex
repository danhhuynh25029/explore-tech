defmodule MyApp.Supervisor do
  use Supervisor

  def start_link do
    Supervisor.start_link(__MODULE__, :ok)
  end

  def init(:ok) do
    children = [
      worker(MyApp.Worker, [])
    ]

    supervise(children, strategy: :one_for_one)
  end
end

defmodule MyApp.Worker do
  def start_link do
    Task.start_link(fn -> loop() end)
  end

  defp loop do
    raise "Oops, something went wrong!"
  end
end

MyApp.Supervisor.start_link()
