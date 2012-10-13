class AddUriToTables < ActiveRecord::Migration
  def self.up
    add_column :columns, :uri, :string
    add_index :columns, :uri
    add_column :rows, :uri, :string
    add_index :rows, :uri
    add_column :grid_mappings, :uri, :string
    add_index :grid_mappings, :uri
    add_column :users, :uri, :string
    add_index :users, :uri
    add_column :roles, :uri, :string
    add_index :roles, :uri
    add_column :uploads, :uri, :string
    add_index :uploads, :uri
    add_column :workspace_sharings, :uri, :string
    add_index :workspace_sharings, :uri
  end

  def self.down
    remove_column :columns, :uri
    remove_column :rows, :uri
    remove_column :grid_mappings, :uri
    remove_column :users, :uri
    remove_column :roles, :uri
    remove_column :uploads, :uri
    remove_column :workspace_sharings, :uri
  end
end
