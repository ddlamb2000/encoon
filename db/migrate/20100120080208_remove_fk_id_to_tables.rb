class RemoveFkIdToTables < ActiveRecord::Migration
  def self.up
    remove_column :grids, :workspace_id
    remove_column :columns, :grid_id
    remove_column :columns, :grid_reference_id
    remove_column :rows, :grid_id

    remove_column :users, :create_user_id
    remove_column :users, :update_user_id
    remove_column :workspaces, :create_user_id
    remove_column :workspaces, :update_user_id
    remove_column :grids, :create_user_id
    remove_column :grids, :update_user_id
    remove_column :columns, :create_user_id
    remove_column :columns, :update_user_id
    remove_column :rows, :create_user_id
    remove_column :rows, :update_user_id

    remove_column :rows, :row_id1
    remove_column :rows, :row_id2
    remove_column :rows, :row_id3
    remove_column :rows, :row_id4
    remove_column :rows, :row_id5
    remove_column :rows, :row_id6
    remove_column :rows, :row_id7
    remove_column :rows, :row_id8
    remove_column :rows, :row_id9
    remove_column :rows, :row_id10
    remove_column :rows, :row_id11
    remove_column :rows, :row_id12
    remove_column :rows, :row_id13
    remove_column :rows, :row_id14
    remove_column :rows, :row_id15
    remove_column :rows, :row_id16
    remove_column :rows, :row_id17
    remove_column :rows, :row_id18
    remove_column :rows, :row_id19
    remove_column :rows, :row_id20
  end

  def self.down
    add_column :grids, :workspace_id, :integer
    add_column :grids, :create_user_id, :integer
    add_column :grids, :update_user_id, :integer

    add_column :columns, :grid_id, :string, :integer
    add_column :columns, :grid_reference_id, :integer
    add_column :columns, :create_user_id, :integer
    add_column :columns, :update_user_id, :integer

    add_column :rows, :grid_id, :integer
    add_column :rows, :create_user_id, :integer
    add_column :rows, :update_user_id, :integer

    add_column :users, :create_user_id, :integer
    add_column :users, :update_user_id, :integer

    add_column :workspaces, :create_user_id, :integer
    add_column :workspaces, :update_user_id, :integer

    add_column :rows, :row_id1, :integer
    add_column :rows, :row_id2, :integer
    add_column :rows, :row_id3, :integer
    add_column :rows, :row_id4, :integer
    add_column :rows, :row_id5, :integer
    add_column :rows, :row_id6, :integer
    add_column :rows, :row_id7, :integer
    add_column :rows, :row_id8, :integer
    add_column :rows, :row_id9, :integer
    add_column :rows, :row_id10, :integer
    add_column :rows, :row_id11, :integer
    add_column :rows, :row_id12, :integer
    add_column :rows, :row_id13, :integer
    add_column :rows, :row_id14, :integer
    add_column :rows, :row_id15, :integer
    add_column :rows, :row_id16, :integer
    add_column :rows, :row_id17, :integer
    add_column :rows, :row_id18, :integer
    add_column :rows, :row_id19, :integer
    add_column :rows, :row_id20, :integer
  end
end
