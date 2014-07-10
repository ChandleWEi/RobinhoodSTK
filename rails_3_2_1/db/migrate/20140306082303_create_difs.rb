class CreateDifs < ActiveRecord::Migration
  def change
    create_table :difs do |t|
      t.integer :stock_map_id
      t.date :record_date
      t.decimal :dif, :precision => 10, :scale => 4
    end
  end
end
