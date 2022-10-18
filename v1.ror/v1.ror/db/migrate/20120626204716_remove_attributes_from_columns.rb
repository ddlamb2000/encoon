class RemoveAttributesFromColumns < ActiveRecord::Migration
  def self.up
    remove_column :columns, :column_reference_uuid
    remove_column :columns, :reference_grid_option_uuid
  end

  def self.down
    add_column :columns, :column_reference_uuid, :string, :limit => 36
    add_column :columns, :reference_grid_option_uuid, :string, :limit => 36
  end
end
