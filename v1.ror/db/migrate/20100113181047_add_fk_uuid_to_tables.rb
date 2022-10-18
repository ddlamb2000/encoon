class AddFkUuidToTables < ActiveRecord::Migration
  def self.up
    add_column :grids, :workspace_uuid, :string, :limit => 36
    add_column :grids, :create_user_uuid, :string, :limit => 36
    add_column :grids, :update_user_uuid, :string, :limit => 36

    add_column :columns, :grid_uuid, :string, :limit => 36
    add_column :columns, :grid_reference_uuid, :string, :limit => 36
    add_column :columns, :create_user_uuid, :string, :limit => 36
    add_column :columns, :update_user_uuid, :string, :limit => 36

    add_column :rows, :grid_uuid, :string, :limit => 36
    add_column :rows, :create_user_uuid, :string, :limit => 36
    add_column :rows, :update_user_uuid, :string, :limit => 36
    add_column :rows, :row_uuid1, :string, :limit => 36
    add_column :rows, :row_uuid2, :string, :limit => 36
    add_column :rows, :row_uuid3, :string, :limit => 36
    add_column :rows, :row_uuid4, :string, :limit => 36
    add_column :rows, :row_uuid5, :string, :limit => 36
    add_column :rows, :row_uuid6, :string, :limit => 36
    add_column :rows, :row_uuid7, :string, :limit => 36
    add_column :rows, :row_uuid8, :string, :limit => 36
    add_column :rows, :row_uuid9, :string, :limit => 36
    add_column :rows, :row_uuid10, :string, :limit => 36
    add_column :rows, :row_uuid11, :string, :limit => 36
    add_column :rows, :row_uuid12, :string, :limit => 36
    add_column :rows, :row_uuid13, :string, :limit => 36
    add_column :rows, :row_uuid14, :string, :limit => 36
    add_column :rows, :row_uuid15, :string, :limit => 36
    add_column :rows, :row_uuid16, :string, :limit => 36
    add_column :rows, :row_uuid17, :string, :limit => 36
    add_column :rows, :row_uuid18, :string, :limit => 36
    add_column :rows, :row_uuid19, :string, :limit => 36
    add_column :rows, :row_uuid20, :string, :limit => 36

    add_column :users, :create_user_uuid, :string, :limit => 36
    add_column :users, :update_user_uuid, :string, :limit => 36

    add_column :workspaces, :create_user_uuid, :string, :limit => 36
    add_column :workspaces, :update_user_uuid, :string, :limit => 36

    add_index :grids, :workspace_uuid
    add_index :columns, :grid_uuid
    add_index :rows, :grid_uuid
  end

  def self.down
    remove_column :grids, :workspace_uuid
    remove_column :columns, :grid_uuid
    remove_column :columns, :grid_reference_uuid
    remove_column :rows, :grid_uuid

    remove_column :users, :create_user_uuid
    remove_column :users, :update_user_uuid
    remove_column :workspaces, :create_user_uuid
    remove_column :workspaces, :update_user_uuid
    remove_column :grids, :create_user_uuid
    remove_column :grids, :update_user_uuid
    remove_column :columns, :create_user_uuid
    remove_column :columns, :update_user_uuid
    remove_column :rows, :create_user_uuid
    remove_column :rows, :update_user_uuid

    remove_column :rows, :row_uuid1
    remove_column :rows, :row_uuid2
    remove_column :rows, :row_uuid3
    remove_column :rows, :row_uuid4
    remove_column :rows, :row_uuid5
    remove_column :rows, :row_uuid6
    remove_column :rows, :row_uuid7
    remove_column :rows, :row_uuid8
    remove_column :rows, :row_uuid9
    remove_column :rows, :row_uuid10
    remove_column :rows, :row_uuid11
    remove_column :rows, :row_uuid12
    remove_column :rows, :row_uuid13
    remove_column :rows, :row_uuid14
    remove_column :rows, :row_uuid15
    remove_column :rows, :row_uuid16
    remove_column :rows, :row_uuid17
    remove_column :rows, :row_uuid18
    remove_column :rows, :row_uuid19
    remove_column :rows, :row_uuid20
  end
end