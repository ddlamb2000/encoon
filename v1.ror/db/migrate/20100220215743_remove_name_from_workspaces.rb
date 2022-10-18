class RemoveNameFromWorkspaces < ActiveRecord::Migration
  def self.up
    remove_column :workspaces, :name
    remove_column :workspaces, :description
  end

  def self.down
    add_column :workspaces, :name, :string
    add_column :workspaces, :description, :string
  end
end
