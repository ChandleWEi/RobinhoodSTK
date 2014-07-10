class Sma < ActiveRecord::Base
  attr_accessible :duration, :record_date, :sma, :stock_map_id
  belongs_to :stock_map
end
