class AddIndexToWorkspaces < ActiveRecord::Migration
  def self.up
    add_index :workspaces, :uri
    add_index :workspaces, :public
  end

  def self.down
    remove_index :workspaces, :uri
    remove_index :workspaces, :public
  end
end