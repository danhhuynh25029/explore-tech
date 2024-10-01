defmodule StockPriceTracker.Application do
  use Application

  def start(_type, _args) do
    children = [
      {StockPriceTracker, []}
    ]

    opts = [strategy: :one_for_one, name: StockPriceTracker.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
