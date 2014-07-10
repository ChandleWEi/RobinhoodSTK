class CreateRawdata < ActiveRecord::Migration
  def change
    create_table :rawdata do |t|
      t.integer :stock_id
      t.date    :record_date
      t.decimal :open_price, precision: 10, scale: 2
      t.decimal :high_price, precision: 10, scale: 2
      t.decimal :low_price, precision: 10, scale: 2
      t.decimal :close_price, precision: 10, scale: 2
      t.integer :volume, limit: 8
      t.decimal :adjusted_close_price, precision: 10, scale: 2

      t.timestamps
    end
  end
end
