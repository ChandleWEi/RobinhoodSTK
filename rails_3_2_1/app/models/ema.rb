class Ema < ActiveRecord::Base
  attr_accessible :duration, :ema, :record_date, :stock_map_id
  belongs_to :stock_map

  #EMA(today) = exp*close_price(today) + (1-exp)*EMA(yesterday) or exp*(close_price(today) - EMA(yesterday))+EMA(yesterday)
  def interate_ema
  end
end
