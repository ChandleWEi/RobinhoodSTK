class Rawdatum < ActiveRecord::Base
  attr_accessible :adjusted_close_price, :close_price, :record_date, :high_price, :low_price, :open_price, :volume, :stock_map_id
  belongs_to :stock_map

  def self.range_records_before_date(stock_map_id, date, duration)
    self.where("stock_map_id = ? and volume <> 0 and record_date <= ?", stock_map_id, date.to_date).order("record_date desc").limit(duration)
  end

  def self.find_data_before_date(stock_map_id, end_date, start_date, order_type)
    self.where("stock_map_id = ? and volume <> 0 and record_date < ? and record_date >= ?", stock_map_id, end_date, start_date).order("record_date#{order_type}")
  end
end
