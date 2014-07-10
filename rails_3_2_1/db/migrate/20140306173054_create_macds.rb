class CreateMacds < ActiveRecord::Migration
  def change
    create_table :macds do |t|
      t.integer :stock_map_id
      t.date :record_date
      t.decimal :macd, :precision => 10, :scale => 4
    end
  end
end
