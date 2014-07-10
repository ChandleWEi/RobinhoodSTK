class Macd < ActiveRecord::Base
  attr_accessible :macd, :record_date, :stock_map_id
  belongs_to :stock_map
end
