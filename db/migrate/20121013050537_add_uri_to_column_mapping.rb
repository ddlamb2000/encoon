class AddUriToColumnMapping < ActiveRecord::Migration
  def self.up
    add_column :column_mappings, :uri, :string
    add_index :column_mappings, :uri
  end

  def self.down
    remove_column :column_mappings, :uri
  end
end