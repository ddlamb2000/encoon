class AddUriToGrid < ActiveRecord::Migration
  def self.up
    add_column :grids, :uri, :string
    add_index :grids, :uri
  end

  def self.down
    remove_column :grids, :uri
  end
end