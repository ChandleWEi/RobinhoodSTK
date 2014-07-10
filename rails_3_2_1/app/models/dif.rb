class Dif < ActiveRecord::Base
  attr_accessible :dif, :record_date, :stock_map_id
  belongs_to :stock_map
end
