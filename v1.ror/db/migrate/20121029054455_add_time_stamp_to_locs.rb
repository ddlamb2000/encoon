class AddTimeStampToLocs < ActiveRecord::Migration
  def self.up
    add_column :grid_locs, :created_at, :datetime
    add_column :grid_locs, :updated_at, :datetime
    add_column :grid_locs, :create_user_uuid, :string, :limit => 36
    add_column :grid_locs, :update_user_uuid, :string, :limit => 36
    add_column :workspace_locs, :created_at, :datetime
    add_column :workspace_locs, :updated_at, :datetime
    add_column :workspace_locs, :create_user_uuid, :string, :limit => 36
    add_column :workspace_locs, :update_user_uuid, :string, :limit => 36
    add_column :column_locs, :created_at, :datetime
    add_column :column_locs, :updated_at, :datetime
    add_column :column_locs, :create_user_uuid, :string, :limit => 36
    add_column :column_locs, :update_user_uuid, :string, :limit => 36
    add_column :row_locs, :created_at, :datetime
    add_column :row_locs, :updated_at, :datetime
    add_column :row_locs, :create_user_uuid, :string, :limit => 36
    add_column :row_locs, :update_user_uuid, :string, :limit => 36
    add_column :role_locs, :created_at, :datetime
    add_column :role_locs, :updated_at, :datetime
    add_column :role_locs, :create_user_uuid, :string, :limit => 36
    add_column :role_locs, :update_user_uuid, :string, :limit => 36
  end

  def self.down
    remove_column :grid_locs, :created_at
    remove_column :grid_locs, :updated_at
    remove_column :grid_locs, :create_user_uuid
    remove_column :grid_locs, :update_user_uuid
    remove_column :workspace_locs, :created_at
    remove_column :workspace_locs, :updated_at
    remove_column :workspace_locs, :create_user_uuid
    remove_column :workspace_locs, :update_user_uuid
    remove_column :column_locs, :created_at
    remove_column :column_locs, :updated_at
    remove_column :column_locs, :create_user_uuid
    remove_column :column_locs, :update_user_uuid
    remove_column :row_locs, :created_at
    remove_column :row_locs, :updated_at
    remove_column :row_locs, :create_user_uuid
    remove_column :row_locs, :update_user_uuid
    remove_column :role_locs, :created_at
    remove_column :role_locs, :updated_at
    remove_column :role_locs, :create_user_uuid
    remove_column :role_locs, :update_user_uuid
  end
end