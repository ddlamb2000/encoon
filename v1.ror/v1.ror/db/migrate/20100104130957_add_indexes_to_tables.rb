class AddIndexesToTables < ActiveRecord::Migration
  def self.up
    add_index :grids, :workspace_id
    add_index :columns, :grid_id
    add_index :rows, :grid_id
  end

  def self.down
    remove_index :grids, :workspace_id
    remove_index :columns, :grid_id
    remove_index :rows, :grid_id
  end
end
