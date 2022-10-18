class AddWhoUpdatedToTables < ActiveRecord::Migration
  def self.up
    add_column :users, :create_user_id, :integer
    add_column :users, :update_user_id, :integer
    add_column :workspaces, :create_user_id, :integer
    add_column :workspaces, :update_user_id, :integer
    add_column :grids, :create_user_id, :integer
    add_column :grids, :update_user_id, :integer
    add_column :columns, :create_user_id, :integer
    add_column :columns, :update_user_id, :integer
  end

  def self.down
    remove_column :users, :create_user_id
    remove_column :users, :update_user_id
    remove_column :workspaces, :create_user_id
    remove_column :workspaces, :update_user_id
    remove_column :grids, :create_user_id
    remove_column :grids, :update_user_id
    remove_column :columns, :create_user_id
    remove_column :columns, :update_user_id
  end
end
