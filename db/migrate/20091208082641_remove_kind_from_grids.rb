class RemoveKindFromGrids < ActiveRecord::Migration
  def self.up
    remove_column :grids, :kind
  end

  def self.down
    add_column :grids, :kind, :string
  end
end
