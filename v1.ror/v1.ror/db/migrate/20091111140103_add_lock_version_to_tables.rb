class AddLockVersionToTables < ActiveRecord::Migration
  def self.up
    add_column :users, :lock_version, :integer, :default => 0
    add_column :workspaces, :lock_version, :integer, :default => 0
    add_column :grids, :lock_version, :integer, :default => 0
    add_column :columns, :lock_version, :integer, :default => 0
  end

  def self.down
    remove_column :users, :lock_version
    remove_column :workspaces, :lock_version
    remove_column :grids, :lock_version
    remove_column :columns, :lock_version
  end
end
