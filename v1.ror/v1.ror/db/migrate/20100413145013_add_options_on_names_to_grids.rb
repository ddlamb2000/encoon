class AddOptionsOnNamesToGrids < ActiveRecord::Migration
  def self.up
    add_column :grids, :has_name, :boolean, :default => true
    add_column :grids, :has_description, :boolean, :default => true
  end

  def self.down
    remove_column :grids, :has_name
    remove_column :grids, :has_description
  end
end
