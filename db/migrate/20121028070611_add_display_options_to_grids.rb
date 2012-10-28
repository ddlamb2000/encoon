class AddDisplayOptionsToGrids < ActiveRecord::Migration
  def self.up
    add_column :grids, :template_uuid, :string, :limit => 36
    add_column :grids, :display_uuid, :string, :limit => 36
    add_column :grids, :sort_uuid, :string, :limit => 36
  end

  def self.down
    remove_column :grids, :template_uuid
    remove_column :grids, :display_uuid
    remove_column :grids, :sort_uuid
  end
end