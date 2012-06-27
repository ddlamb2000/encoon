class AddSecurityOptionsToWorkspaces < ActiveRecord::Migration
  def self.up
    add_column :workspaces, :public, :boolean
    add_column :workspaces, :default_role_uuid, :string, :limit => 36
    add_column :workspaces, :uri, :string
  end

  def self.down
    remove_column :workspaces, :public
    remove_column :workspaces, :default_role_uuid
    remove_column :workspaces, :uri
  end
end
