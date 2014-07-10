class CreateSmas < ActiveRecord::Migration
  def change
    create_table :smas do |t|
      t.integer :stock_map_id
      t.date :record_date
      t.integer :duration
      t.decimal :sma, :precision => 10, :scale => 4
    end
  end
end
