defmodule StockPriceTracker do
  use GenServer

  def start_link do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def init(state) do
    {:ok, state}
  end

  def handle_info(:update_stock_price, state) do
    # Mô phỏng việc cập nhật giá cổ phiếu
    new_price = 100 + Enum.random(1..10)
    IO.puts "New stock price: $#{new_price}"
    {:noreply, state}
  end
end
