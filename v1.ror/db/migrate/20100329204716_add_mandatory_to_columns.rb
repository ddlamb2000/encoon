class AddMandatoryToColumns < ActiveRecord::Migration
  def self.up
    add_column :columns, :required, :boolean
    add_column :columns, :length, :integer
    add_column :columns, :precision, :integer
    add_column :columns, :regex, :string
    add_column :columns, :column_reference_uuid, :string, :limit => 36
    add_column :columns, :reference_grid_option_uuid, :string, :limit => 36
  end

  def self.down
    remove_column :columns, :required
    remove_column :columns, :length
    remove_column :columns, :precision
    remove_column :columns, :regex
    remove_column :columns, :column_reference_uuid
    remove_column :columns, :reference_grid_option_uuid
  end
end
