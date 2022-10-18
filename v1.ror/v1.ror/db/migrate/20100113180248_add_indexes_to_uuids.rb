class AddIndexesToUuids < ActiveRecord::Migration
  def self.up
    add_index :users, :uuid
    add_index :workspaces, :uuid
    add_index :grids, :uuid
    add_index :columns, :uuid
    add_index :rows, :uuid
  end

  def self.down
    remove_index :users, :uuid
    remove_index :workspaces, :uuid
    remove_index :grids, :uuid
    remove_index :columns, :uuid
    remove_index :rows, :uuid
  end
end
