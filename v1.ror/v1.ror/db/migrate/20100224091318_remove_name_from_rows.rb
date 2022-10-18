class RemoveNameFromRows < ActiveRecord::Migration
  def self.up
    remove_column :rows, :name
    remove_column :rows, :description
  end

  def self.down
    add_column :rows, :name, :string
    add_column :rows, :description, :string
  end
end
