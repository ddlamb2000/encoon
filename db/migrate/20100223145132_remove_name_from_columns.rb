class RemoveNameFromColumns < ActiveRecord::Migration
  def self.up
    remove_column :columns, :name
    remove_column :columns, :description
  end

  def self.down
    add_column :columns, :name, :string
    add_column :columns, :description, :string
  end
end
