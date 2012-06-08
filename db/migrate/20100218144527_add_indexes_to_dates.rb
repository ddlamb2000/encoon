class AddIndexesToDates < ActiveRecord::Migration
  def self.up
    add_index :users, [:uuid, :begin, :end]
    add_index :workspaces, [:uuid, :begin, :end]
    add_index :grids, [:uuid, :begin, :end]
    add_index :columns, [:uuid, :begin, :end]
    add_index :rows, [:uuid, :begin, :end]
  end

  def self.down
    remove_index :users, [:uuid, :begin, :end]
    remove_index :workspaces, [:uuid, :begin, :end]
    remove_index :grids, [:uuid, :begin, :end]
    remove_index :columns, [:uuid, :begin, :end]
    remove_index :rows, [:uuid, :begin, :end]
  end
end
