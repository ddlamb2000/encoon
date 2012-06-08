class RemoveNameFromGrids < ActiveRecord::Migration
  def self.up
    remove_column :grids, :name
    remove_column :grids, :description
  end

  def self.down
    add_column :grids, :name, :string
    add_column :grids, :description, :string
  end
end
