class RemoveRevisionFromTables < ActiveRecord::Migration
  def self.up
    remove_column :workspaces, :revision
    remove_column :grids, :revision
    remove_column :columns, :revision
  end

  def self.down
    add_column :workspaces, :revision, :integer
    add_column :grids, :revision, :integer
    add_column :columns, :revision, :integer
  end
end
