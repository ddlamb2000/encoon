class AddGridToColumn < ActiveRecord::Migration
  def self.up
    add_column :columns, :grid_reference_id, :integer
  end

  def self.down
    remove_column :columns, :grid_reference_id
  end
end
