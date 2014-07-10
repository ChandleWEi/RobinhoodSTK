class CreateEmas < ActiveRecord::Migration
  def change
    create_table :emas do |t|
      t.integer :stock_map_id
      t.date :record_date
      t.integer :duration
      t.decimal :ema, :precision => 10, :scale => 4
    end
  end
end
