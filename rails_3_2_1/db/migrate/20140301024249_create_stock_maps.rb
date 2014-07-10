class CreateStockMaps < ActiveRecord::Migration
  def change
    create_table :stock_maps do |t|
      t.string :stock_id
      t.string :stock_name

    end
  end
end
